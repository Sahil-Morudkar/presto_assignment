package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Sahil-Morudkar/presto_assignment/internal/handler"
)

func NewRouter(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	pricingHandler := handler.NewPricingHandler(db)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/chargers/{chargerId}/pricing", pricingHandler.GetPricing)
	})

	return r
}