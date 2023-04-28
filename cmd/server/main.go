package main

import (
	"net/http"

	"github.com/Zagir2000/alert/cmd/server/parser"
	"github.com/Zagir2000/alert/cmd/server/storage"
)

func CollectMetricsAndALerts(res http.ResponseWriter, req *http.Request) {
	var storage storage.MemStorageUsage = &storage.MemStorage{}
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)

	}
	err := storage.CollectMetricsAndALerts(req.URL.String())
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
	return

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, CollectMetricsAndALerts)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
