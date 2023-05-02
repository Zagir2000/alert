package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zagir2000/alert/internal/server/handlers"
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
)

func run(r *chi.Mux) error {
	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, r)
}

func main() {
	m := storage.NewMemStorage()
	NewHandStruct := handlers.MetricHandlerNew(m)

	r := chi.NewRouter()
	// r.Get("/value/*", handlers.GetMetric)
	// r.Get("/", handlers.ShowMetrics)
	r.Post("/update/{metricType}/{metricName}/{value}", NewHandStruct.CollectMetricsAndALerts)
	r.Get("/", NewHandStruct.MainPage)
	r.Get("/value/{metricType}/{metricName}", NewHandStruct.NowValueMetrics)
	parseFlags()
	if err := run(r); err != nil {
		log.Fatalln(err)
	}

}
