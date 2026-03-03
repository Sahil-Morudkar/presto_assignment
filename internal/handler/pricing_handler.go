package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Sahil-Morudkar/presto_assignment/internal/model"
	"github.com/Sahil-Morudkar/presto_assignment/internal/service"
)

type PricingHandler struct {
	service *service.PricingService
}

func NewPricingHandler(service *service.PricingService) *PricingHandler {
	return &PricingHandler{service: service}
}

func (h *PricingHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {

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

	var req model.CreateScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	err := h.service.CreatePricingSchedule(r.Context(), chargerID, req)
	if err != nil {

		// You can improve this by mapping specific errors
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

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

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: "date query parameter is required (YYYY-MM-DD)",
		})
		return
	}

	result, err := h.service.GetDailyPricing(r.Context(), chargerID, dateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.APIResponse{
		Status:  "success",
		Message: "pricing retrieved successfully",
		Data:    result,
	})
}