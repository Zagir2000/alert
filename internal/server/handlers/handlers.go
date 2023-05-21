package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/models"
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MetricHandler struct {
	Storage storage.Repository
}

func MetricHandlerNew(s storage.Repository) *MetricHandler {
	return &MetricHandler{Storage: s}
}

func (m *MetricHandler) AllMetrics() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Debug("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		res.Write([]byte("<h1>Gauge metrics</h1>"))
		AllGaugeValues := m.Storage.GetAllGaugeValues()
		OrderAllGaugeValues := make([]string, 0, len(AllGaugeValues))
		for k := range AllGaugeValues {
			OrderAllGaugeValues = append(OrderAllGaugeValues, k)
		}
		// sort the slice by keys
		sort.Strings(OrderAllGaugeValues)
		for _, k := range OrderAllGaugeValues {
			fmt.Fprintf(res, "%s: %g\n", k, AllGaugeValues[k])
		}
		res.Write([]byte("<h1>Counter metrics</h1>"))
		AllCounterValues := m.Storage.GetAllCounterValues()
		for k, v := range AllCounterValues {
			fmt.Fprintf(res, "%s: %d\n", k, v)
		}
		res.Header().Add("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
	}
}

func (m *MetricHandler) NowValueMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", req.Method))
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
			return
		}
		res.Write([]byte(fmt.Sprintf("%d", value)))
	case "gauge":
		value, ok := m.Storage.GetGauge(metricName)
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.Write([]byte(fmt.Sprintf("%g", value)))
	default:
		{
			logger.Log.Debug("could not determine the type of metric")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	res.WriteHeader(http.StatusOK)
}

func (m *MetricHandler) NewMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", req.Method))
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
			logger.Log.Debug("cannot parse to int64", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		err = m.Storage.AddCounterValue(metricName, valueint64)
		if err != nil {
			logger.Log.Debug("cannot get new counter value", zap.Error(err))
			res.WriteHeader(http.StatusNotFound)
			return
		}
	case "gauge":
		valuefloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Log.Debug("cannot parse to float64", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		err = m.Storage.AddGaugeValue(metricName, valuefloat64)
		if err != nil {
			logger.Log.Debug("cannot add new gauge value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
	default:
		{
			logger.Log.Debug("could not determine the type of metric")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

func (m *MetricHandler) NowValueMetricsToJson(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", req.Method))
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var jsonMetrics models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&jsonMetrics); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// заполняем модель ответа
	switch strings.ToLower(jsonMetrics.MType) {
	case "counter":
		value, ok := m.Storage.GetCounter(jsonMetrics.ID)
		if !ok {
			logger.Log.Debug("cannot get counter value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		jsonMetrics.Delta = &value
	case "gauge":
		value, ok := m.Storage.GetGauge(jsonMetrics.ID)
		if !ok {
			logger.Log.Debug("cannot get gauge value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		jsonMetrics.Value = &value
	default:
		{
			logger.Log.Debug("could not determine the type of metric")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	response, err := json.Marshal(jsonMetrics)
	if err != nil {
		logger.Log.Debug("cannot marshal to json", zap.Error(err))
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (m *MetricHandler) NewMetricsToJson(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", req.Method))
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var jsonMetrics models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&jsonMetrics); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// заполняем модель ответа
	switch strings.ToLower(jsonMetrics.MType) {
	case "counter":
		err := m.Storage.AddCounterValue(jsonMetrics.ID, *jsonMetrics.Delta)
		if err != nil {
			logger.Log.Debug("cannot add new counter value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		value, ok := m.Storage.GetCounter(jsonMetrics.ID)
		if !ok {
			logger.Log.Debug("cannot get counter value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		jsonMetrics.Delta = &value
	case "gauge":
		err := m.Storage.AddGaugeValue(jsonMetrics.ID, *jsonMetrics.Value)
		if err != nil {
			logger.Log.Debug("cannot add new gauge value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		value, ok := m.Storage.GetGauge(jsonMetrics.ID)
		if !ok {
			logger.Log.Debug("cannot get gauge value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		jsonMetrics.Value = &value
	default:
		{
			logger.Log.Debug("could not determine the type of metric")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	response, err := json.Marshal(jsonMetrics)
	if err != nil {
		logger.Log.Debug("cannot marshal to json", zap.Error(err))
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}
