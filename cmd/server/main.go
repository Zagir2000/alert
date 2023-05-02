package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Zagir2000/alert/internal/server/handlers"
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
)

var flagRunAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
}

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
