package handlers

import (
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	m := storage.NewMemStorage()
	NewHandStruct := MetricHandlerNew(m)
	r := chi.NewRouter()
	// r.Get("/value/*", handlers.GetMetric)
	// r.Get("/", handlers.ShowMetrics)
	r.Post("/update/{metricType}/{metricName}/{value}", NewHandStruct.NewMetrics)
	r.Get("/", NewHandStruct.AllMetrics)
	r.Get("/value/{metricType}/{metricName}", NewHandStruct.NowValueMetrics)
	return r
}
