package handlers

import (
	"github.com/Zagir2000/alert/internal/logger"
	"github.com/go-chi/chi/v5"
)

func Router(newHandStruct *MetricHandler) chi.Router {

	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Post("/update/", gzipMiddleware(newHandStruct.NewMetricsToJSON))
	r.Post("/value/", gzipMiddleware(newHandStruct.AddValueMetricsToJSON))
	r.Post("/update/{metricType}/{metricName}/{value}", gzipMiddleware(newHandStruct.UpdateNewMetrics))
	r.Get("/", gzipMiddleware(newHandStruct.GetAllMetrics))
	r.Get("/value/{metricType}/{metricName}", gzipMiddleware(newHandStruct.GetNowValueMetrics))

	return r
}
