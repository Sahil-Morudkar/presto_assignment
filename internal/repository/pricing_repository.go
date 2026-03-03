package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Sahil-Morudkar/presto_assignment/internal/model"
)

type PricingRepository struct {
	DB *pgxpool.Pool
}

func NewPricingRepository(db *pgxpool.Pool) *PricingRepository {
	return &PricingRepository{DB: db}
}

// ChargerExists checks if a charger exists
func (r *PricingRepository) ChargerExists(ctx context.Context, chargerID string) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM chargers WHERE id = $1
		)
	`, chargerID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PricingRepository) CreateSchedule(
	ctx context.Context,
	chargerID string,
	effectiveFrom time.Time,
	periods []model.PricingPeriodInput,
) error {

	// Begin transaction to ensure atomicity
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Close existing active schedule
	_, err = tx.Exec(ctx, `
		UPDATE tou_pricing_schedules
		SET effective_to = $1
		WHERE charger_id = $2
		AND effective_to IS NULL
	`, effectiveFrom.AddDate(0, 0, -1), chargerID)

	if err != nil {
		return err
	}

	// Insert new schedule
	var scheduleID string
	err = tx.QueryRow(ctx, `
		INSERT INTO tou_pricing_schedules (charger_id, effective_from)
		VALUES ($1, $2)
		RETURNING id
	`, chargerID, effectiveFrom).Scan(&scheduleID)

	if err != nil {
		return err
	}

	// Insert associated pricing periods
	for _, p := range periods {
		_, err = tx.Exec(ctx, `
			INSERT INTO tou_pricing_periods
			(schedule_id, start_time, end_time, price_per_kwh)
			VALUES ($1, $2, $3, $4)
		`, scheduleID, p.StartTime, p.EndTime, p.PricePerKwh)

		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit(ctx)
}

// GetDailyPricing fetches full schedule for a given charger and date
func (r *PricingRepository) GetDailyPricing(
	ctx context.Context,
	chargerID string,
	inputDate time.Time,
) (*model.DailyPricingResponse, error) {

	// Step 1: Fetch applicable schedule for the date
	var scheduleID string
	var effectiveFrom time.Time

	err := r.DB.QueryRow(ctx, `
		SELECT id, effective_from
		FROM tou_pricing_schedules
		WHERE charger_id = $1
		  AND effective_from <= $2
		  AND (effective_to IS NULL OR effective_to >= $2)
		ORDER BY effective_from DESC
		LIMIT 1
	`, chargerID, inputDate).Scan(&scheduleID, &effectiveFrom)

	if err != nil {
		return nil, err
	}

	// Step 2: Fetch all pricing periods for that schedule
	rows, err := r.DB.Query(ctx, `
		SELECT start_time, end_time, price_per_kwh
		FROM tou_pricing_periods
		WHERE schedule_id = $1
		ORDER BY start_time ASC
	`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []model.PricingPeriodResponse

	for rows.Next() {

		var startTime time.Time
		var endTime time.Time
		var price float64

		if err := rows.Scan(&startTime, &endTime, &price); err != nil {
			return nil, err
		}

		periods = append(periods, model.PricingPeriodResponse{
			StartTime:   startTime.Format("15:04"),
			EndTime:     endTime.Format("15:04"),
			PricePerKwh: price,
		})
	}

	return &model.DailyPricingResponse{
		ChargerID:     chargerID,
		EffectiveFrom: effectiveFrom.Format("2006-01-02"),
		Periods:       periods,
	}, nil
}

// CreateBulkSchedule applies same schedule to multiple chargers atomically.
// If any charger fails, entire transaction is rolled back.
// Returns failed charger ID (if any).
func (r *PricingRepository) CreateBulkSchedule(
	ctx context.Context,
	chargerIDs []string,
	effectiveFrom time.Time,
	periods []model.PricingPeriodInput,
) (string, error) {

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	for _, chargerID := range chargerIDs {

		// Close previous active schedule
		_, err := tx.Exec(ctx, `
			UPDATE tou_pricing_schedules
			SET effective_to = $1
			WHERE charger_id = $2
			AND effective_to IS NULL
		`, effectiveFrom.AddDate(0, 0, -1), chargerID)

		if err != nil {
			return chargerID, err
		}

		// Insert new schedule
		var scheduleID string
		err = tx.QueryRow(ctx, `
			INSERT INTO tou_pricing_schedules (charger_id, effective_from)
			VALUES ($1, $2)
			RETURNING id
		`, chargerID, effectiveFrom).Scan(&scheduleID)

		if err != nil {
			return chargerID, err
		}

		// Insert pricing periods
		for _, p := range periods {
			_, err = tx.Exec(ctx, `
				INSERT INTO tou_pricing_periods
				(schedule_id, start_time, end_time, price_per_kwh)
				VALUES ($1, $2, $3, $4)
			`, scheduleID, p.StartTime, p.EndTime, p.PricePerKwh)

			if err != nil {
				return chargerID, err
			}
		}
	}

	// Commit only if all succeed
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return "", nil
}