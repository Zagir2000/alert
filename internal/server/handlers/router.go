package handlers

import (
	"context"
	"net/http/pprof"
	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/Zagir2000/alert/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Роутер, который реализует все методы API.
func Router(ctx context.Context, log *zap.Logger, newHandStruct *MetricHandlerDB, keySecret string) chi.Router {
	r := chi.NewRouter()
	r.Use((logger.WithLogging(log)))
	r.Post("/update/", gzipMiddleware(newHandStruct.NewMetricsToJSON(ctx, log)))
	r.Post("/updates/", gzipMiddleware(newHandStruct.UpdateNewMetricsBatch(ctx, log, keySecret)))
	r.Post("/value/", gzipMiddleware(newHandStruct.AddValueMetricsToJSON(ctx, log)))
	r.Post("/update/{metricType}/{metricName}/{value}", gzipMiddleware(newHandStruct.UpdateNewMetrics(ctx, log)))
	r.Get("/", gzipMiddleware(newHandStruct.GetAllMetrics(ctx, log)))
	r.Get("/value/{metricType}/{metricName}", gzipMiddleware(newHandStruct.GetNowValueMetrics(ctx, log)))
	r.Get("/ping", gzipMiddleware(newHandStruct.PingDBConnect(ctx, log)))

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	return r
}
