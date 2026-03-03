package model

// PricingPeriodInput represents one pricing time bucket
// Example: 00:00 → 06:00 @ 0.15
type PricingPeriodInput struct {
	StartTime   string  `json:"start_time"`     // Raw input time (will be normalized)
	EndTime     string  `json:"end_time"`       // Raw input time (will be normalized)
	PricePerKwh float64 `json:"price_per_kwh"`  // Price in $/kWh
}

// CreateScheduleRequest represents payload to create a schedule
type CreateScheduleRequest struct {
	EffectiveFrom string               `json:"effective_from"` // Expected format: YYYY-MM-DD
	Periods       []PricingPeriodInput `json:"periods"`        // List of pricing buckets
}

// PricingPeriodResponse represents a single pricing period returned to client
type PricingPeriodResponse struct {
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	PricePerKwh float64 `json:"price_per_kwh"`
}

// DailyPricingResponse represents full schedule for a given date
type DailyPricingResponse struct {
	ChargerID     string                  `json:"charger_id"`
	EffectiveFrom string                  `json:"effective_from"`
	Periods       []PricingPeriodResponse `json:"periods"`
}

// BulkCreateScheduleRequest represents bulk pricing update request
type BulkCreateScheduleRequest struct {
	ChargerIDs    []string             `json:"charger_ids"`
	EffectiveFrom string               `json:"effective_from"`
	Periods       []PricingPeriodInput `json:"periods"`
}