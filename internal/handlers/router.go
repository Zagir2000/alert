package handlers

import (
	"context"

	"github.com/Zagir2000/alert/internal/logger"
	"github.com/go-chi/chi/v5"
)

func Router(ctx context.Context, newHandStruct *MetricHandlerDB) chi.Router {
	zapNewLogger := logger.NewZapLoggerStruct(newHandStruct.log)
	r := chi.NewRouter()
	r.Use(zapNewLogger.WithLogging)
	r.Post("/update/", gzipMiddleware(newHandStruct.NewMetricsToJSON(ctx)))
	r.Post("/updates/", gzipMiddleware(newHandStruct.UpdateNewMetricsBatch(ctx)))
	r.Post("/value/", gzipMiddleware(newHandStruct.AddValueMetricsToJSON(ctx)))
	r.Post("/update/{metricType}/{metricName}/{value}", gzipMiddleware(newHandStruct.UpdateNewMetrics(ctx)))
	r.Get("/", gzipMiddleware(newHandStruct.GetAllMetrics(ctx)))
	r.Get("/value/{metricType}/{metricName}", gzipMiddleware(newHandStruct.GetNowValueMetrics(ctx)))
	r.Get("/ping", gzipMiddleware(newHandStruct.PingDBConnect(ctx)))
	return r
}
