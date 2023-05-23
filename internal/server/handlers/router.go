package handlers

import (
	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	m := storage.NewMemStorage()
	newHandStruct := MetricHandlerNew(m)
	r := chi.NewRouter()
	// r.Get("/value/*", handlers.GetMetric)
	// r.Get("/", handlers.ShowMetrics)
	r.Use(logger.WithLogging)
	r.Post("/update/", newHandStruct.NewMetricsToJSON())
	r.Post("/value/", newHandStruct.NowValueMetricsToJSON())
	r.Post("/update/{metricType}/{metricName}/{value}", newHandStruct.NewMetrics)
	r.Get("/", newHandStruct.AllMetrics())
	r.Get("/value/{metricType}/{metricName}", newHandStruct.NowValueMetrics)

	return r
}
