package main

import (
	"log"
	"net/http"

	"github.com/Zagir2000/alert/internal/server/handlers"
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	m := storage.NewMemStorage()
	NewHandStruct := handlers.MetricHandler{m}

	r := chi.NewRouter()
	// r.Get("/value/*", handlers.GetMetric)
	// r.Get("/", handlers.ShowMetrics)
	r.Post("/update/{metricType}/{metricName}/{value}", NewHandStruct.CollectMetricsAndALerts)
	r.Get("/", NewHandStruct.MainPage)
	r.Get("/value/{metricType}/{metricName}", NewHandStruct.NowValueMetrics)
	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		log.Fatalln(err)
	}
}
