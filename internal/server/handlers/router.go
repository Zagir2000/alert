package handlers

import (
	"context"

	"github.com/Zagir2000/alert/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Router(ctx context.Context, log *zap.Logger, newHandStruct *MetricHandlerDB) chi.Router {
	zapNewLogger := logger.NewZapLoggerStruct(log)
	r := chi.NewRouter()
	r.Use(zapNewLogger.WithLogging)
	r.Post("/update/", gzipMiddleware(newHandStruct.NewMetricsToJSON(ctx, log)))
	r.Post("/updates/", gzipMiddleware(newHandStruct.UpdateNewMetricsBatch(ctx, log)))
	r.Post("/value/", gzipMiddleware(newHandStruct.AddValueMetricsToJSON(ctx, log)))
	r.Post("/update/{metricType}/{metricName}/{value}", gzipMiddleware(newHandStruct.UpdateNewMetrics(ctx, log)))
	r.Get("/", gzipMiddleware(newHandStruct.GetAllMetrics(ctx, log)))
	r.Get("/value/{metricType}/{metricName}", gzipMiddleware(newHandStruct.GetNowValueMetrics(ctx, log)))
	r.Get("/ping", gzipMiddleware(newHandStruct.PingDBConnect(ctx, log)))
	return r
}
