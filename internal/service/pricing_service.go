package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Sahil-Morudkar/presto_assignment/internal/model"
	"github.com/Sahil-Morudkar/presto_assignment/internal/repository"
)

type PricingService struct {
	repo *repository.PricingRepository
}

func NewPricingService(repo *repository.PricingRepository) *PricingService {
	return &PricingService{repo: repo}
}

// CreatePricingSchedule validates business rules before saving
func (s *PricingService) CreatePricingSchedule(
	ctx context.Context,
	chargerID string,
	req model.CreateScheduleRequest,
) error {

	// Check charger existence
	exists, err := s.repo.ChargerExists(ctx, chargerID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("charger not found")
	}

	// Validate effective_from format (YYYY-MM-DD)
	effectiveFrom, err := time.Parse(time.DateOnly, req.EffectiveFrom)
	if err != nil {
		return errors.New("effective_from must be in YYYY-MM-DD format")
	}

	// Validate periods exist
	if len(req.Periods) == 0 {
		return errors.New("at least one pricing period is required")
	}

	// Normalize + validate time buckets
	for i := range req.Periods {

		// Normalize start time
		startStr, err := normalizeTime(req.Periods[i].StartTime)
		if err != nil {
			return err
		}

		// Normalize end time
		endStr, err := normalizeTime(req.Periods[i].EndTime)
		if err != nil {
			return err
		}

		req.Periods[i].StartTime = startStr
		req.Periods[i].EndTime = endStr

		// Convert to time.Time for comparison
		startTime, _ := time.Parse("15:04", startStr)
		endTime, _ := time.Parse("15:04", endStr)

		// Validate start < end (no cross-midnight allowed)
		if !startTime.Before(endTime) {
			return errors.New("start_time must be before end_time")
		}

		// Validate price
		if req.Periods[i].PricePerKwh < 0 {
			return errors.New("price_per_kwh must be >= 0")
		}
	}

	// Sort periods by start time
	sort.Slice(req.Periods, func(i, j int) bool {
		return req.Periods[i].StartTime < req.Periods[j].StartTime
	})

	// Check for overlapping periods
	for i := 1; i < len(req.Periods); i++ {
		prevEnd := req.Periods[i-1].EndTime
		currStart := req.Periods[i].StartTime

		if currStart < prevEnd {
			return errors.New("pricing periods cannot overlap")
		}
	}

	// Delegate to repository
	return s.repo.CreateSchedule(ctx, chargerID, effectiveFrom, req.Periods)
}

// GetDailyPricing validates input and delegates to repository
func (s *PricingService) GetDailyPricing(
	ctx context.Context,
	chargerID string,
	dateStr string,
) (*model.DailyPricingResponse, error) {

	// Validate charger exists
	exists, err := s.repo.ChargerExists(ctx, chargerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("charger not found")
	}

	// Parse date (YYYY-MM-DD format expected)
	inputDate, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return nil, errors.New("date must be in YYYY-MM-DD format")
	}

	return s.repo.GetDailyPricing(ctx, chargerID, inputDate)
}