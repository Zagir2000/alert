package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
)

type MetricHandler struct {
	Storage storage.Repository
}

func MetricHandlerNew(s storage.Repository) *MetricHandler {
	return &MetricHandler{Storage: s}
}

func (m *MetricHandler) MainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res.Header().Add("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("<h1>Gauge metrics</h1>"))
	AllGaugeValues := m.Storage.GetAllGaugeValues()
	for k, v := range AllGaugeValues {
		res.Write([]byte(fmt.Sprintf("%s: %g\n", k, v)))
	}
	res.Write([]byte("<h1>Counter metrics</h1>"))
	AllCounterValues := m.Storage.GetAllCounterValues()
	for k, v := range AllCounterValues {
		res.Write([]byte(fmt.Sprintf("%s: %d\n", k, v)))
	}
	res.WriteHeader(http.StatusOK)
}

func (m *MetricHandler) NowValueMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	switch metricType {
	case "counter":
		value, ok := m.Storage.GetCounter(metricName)
		if !ok {
			res.WriteHeader(http.StatusNotFound)
		}
		res.Write([]byte(fmt.Sprintf("%d", value)))
	case "gauge":
		value, ok := m.Storage.GetGauge(metricName)
		if !ok {
			res.WriteHeader(http.StatusNotFound)
		}
		res.Write([]byte(fmt.Sprintf("%g", value)))
	default:
		{
			res.WriteHeader(http.StatusBadRequest)
		}
	}
	res.WriteHeader(http.StatusOK)
}

func (m *MetricHandler) Update(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	value := chi.URLParam(req, "value")
	//prepare metric and set value
	switch metricType {
	case "counter":
		valueint64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		}
		m.Storage.AddCounterValue(metricName, valueint64)
	case "gauge":
		valuefloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		}
		m.Storage.AddGaugeValue(metricName, valuefloat64)
	default:
		{
			res.WriteHeader(http.StatusBadRequest)
		}
	}
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

// func MetricHandlerNew(s storage.Repository) *MetricHandler {
// 	return &MetricHandler{Storage: s}
// }
