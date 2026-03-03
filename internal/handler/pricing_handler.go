package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Sahil-Morudkar/presto_assignment/internal/model"
	"github.com/Sahil-Morudkar/presto_assignment/internal/service"
)


// PricingHandler handles HTTP requests related to pricing schedules
type PricingHandler struct {
	service *service.PricingService
}


// NewPricingHandler creates a new instance of PricingHandler with the given service
func NewPricingHandler(service *service.PricingService) *PricingHandler {
	return &PricingHandler{service: service}
}

// CreateSchedule handles POST /chargers/{chargerID}/pricing-schedules
func (h *PricingHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Extract chargerID from URL path
	chargerID := chi.URLParam(r, "chargerID")
	if chargerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: "chargerID is required",
		})
		return
	}

	var req model.CreateScheduleRequest

	// Decode JSON body into struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// Respond with error if JSON is invalid
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	err := h.service.CreatePricingSchedule(r.Context(), chargerID, req)
	if err != nil {
		// Respond with error if service returns an error (e.g., validation failure)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Respond with success message if schedule creation is successful
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.APIResponse{
		Status:  "success",
		Message: "Pricing schedule created successfully",
		Data:    nil,
	})
}

// GetDailyPricing handles GET /chargers/{chargerID}/pricing?date=YYYY-MM-DD
func (h *PricingHandler) GetDailyPricing(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	chargerID := chi.URLParam(r, "chargerID")
	if chargerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: "chargerID is required",
		})
		return
	}

	//
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		//
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: "date query parameter is required (YYYY-MM-DD)",
		})
		return
	}

	// Delegate to service layer to fetch pricing for the given charger and date
	result, err := h.service.GetDailyPricing(r.Context(), chargerID, dateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// Respond with error message from service (e.g., charger not found, invalid date format)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Respond with the retrieved pricing data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.APIResponse{
		Status:  "success",
		Message: "pricing retrieved successfully",
		Data:    result,
	})
}

// CreateBulkSchedule handles POST /pricing-schedules/bulk
func (h *PricingHandler) CreateBulkSchedule(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req model.BulkCreateScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	//
	err := h.service.CreateBulkPricingSchedule(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.APIResponse{
		Status:  "success",
		Message: "bulk pricing schedule created successfully",
	})
}