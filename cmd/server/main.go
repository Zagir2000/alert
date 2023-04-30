package main

import (
	"log"
	"net/http"

	"github.com/Zagir2000/alert/internal/parser"
	"github.com/Zagir2000/alert/internal/storage"
)

func CollectMetricsAndALerts(res http.ResponseWriter, req *http.Request) {
	var storage storage.Repository = &storage.MemStorage{Gaugedata: make(map[string]float64), Counterdata: make(map[string]int64)}
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)

	}
	err := storage.CollectMetricsAndALerts(req.RequestURI)
	if err != nil {
		switch err {
		case parser.ErrType:
			res.WriteHeader(http.StatusBadRequest)
		case parser.ErrValue:
			res.WriteHeader(http.StatusBadRequest)
		case parser.ErrNameMetric:
			res.WriteHeader(http.StatusNotFound)
		}
		return
	}
	res.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, CollectMetricsAndALerts)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		log.Fatalln(err)
	}
}
