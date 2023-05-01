package handlers

import (
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

func (m *MetricHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {

}
func (m *MetricHandler) CollectMetricsAndALerts(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
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
		m.Storage.SetCounter(metricName, valueint64)
	case "gauge":
		valuefloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		}
		m.Storage.SetGauge(metricName, valuefloat64)
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
