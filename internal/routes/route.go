package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Sahil-Morudkar/presto_assignment/internal/handler"
	"github.com/Sahil-Morudkar/presto_assignment/internal/repository"
	"github.com/Sahil-Morudkar/presto_assignment/internal/service"
)

func NewRouter(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	// Layer wiring
	repo := repository.NewPricingRepository(db)
	svc := service.NewPricingService(repo)
	pricingHandler := handler.NewPricingHandler(svc)

	// Route
	r.Post("/chargers/{chargerID}/pricing-schedules", pricingHandler.CreateSchedule)
	r.Get("/chargers/{chargerID}/pricing", pricingHandler.GetDailyPricing)
	r.Post("/pricing-schedules/bulk", pricingHandler.CreateBulkSchedule)

	return r
}