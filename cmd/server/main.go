package main

import (
	"errors"
	"net/http"
	"strconv"
)

type MemStorageUsage interface {
	CollectMetricsAndALerts(res http.ResponseWriter, req *http.Request)
}

func (c MemStorage) CollectMetricsAndALerts(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		// разрешаем только POST-запросы
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	typemetric, value := parser.parseuri(res, req, req.RequestURI)
	if typemetric == "" && typemetric == "" {
		return
	}

	if typemetric == "counter" {
		valueint64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if valueint64 < 0 {
			panic(errors.New("gauge cannot decrease in value"))
		}
		c.counterdata = map[string]counter{"counter": counter(valueint64) + c.counterdata["counter"]}
	} else {
		valuefloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if valuefloat64 < 0 {
			panic(errors.New("counter cannot decrease in value"))
		}
		c.gaugedata = map[string]gauge{"gauge": gauge(valuefloat64) + c.gaugedata["gauge"]}
	}

	res.WriteHeader(http.StatusOK)
}

func main() {

	var data MemStorageUsage = MemStorage{}
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, data.CollectMetricsAndALerts)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
