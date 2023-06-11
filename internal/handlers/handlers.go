package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Zagir2000/alert/internal/models"
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MetricHandlerDB struct {
	Storage storage.Repository
	pgDB    *storage.PostgresDB
	log     *zap.Logger
}

func MetricHandlerNew(s storage.Repository, logger *zap.Logger, pgDB *storage.PostgresDB) *MetricHandlerDB {
	return &MetricHandlerDB{
		Storage: s,
		pgDB:    pgDB,
		log:     logger,
	}
}

func (m *MetricHandlerDB) GetAllMetrics(ctx context.Context) http.HandlerFunc {
	getAllFn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			m.log.Debug("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		res.Header().Add("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		AllGaugeValues := m.Storage.GetAllGaugeValues(ctx)
		OrderAllGaugeValues := make([]string, 0, len(AllGaugeValues))
		for k := range AllGaugeValues {
			OrderAllGaugeValues = append(OrderAllGaugeValues, k)
		}
		// sort the slice by keys
		sort.Strings(OrderAllGaugeValues)
		for _, k := range OrderAllGaugeValues {
			fmt.Fprintf(res, "%s: %g\n", k, AllGaugeValues[k])
		}

		AllCounterValues := m.Storage.GetAllCounterValues(ctx)
		for k, v := range AllCounterValues {
			fmt.Fprintf(res, "%s: %d\n", k, v)
		}

	}
	return getAllFn
}

func (m *MetricHandlerDB) GetNowValueMetrics(ctx context.Context) http.HandlerFunc {
	getValueFn := func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			m.log.Debug("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")

		switch metricType {
		case "counter":
			value, ok := m.Storage.GetCounter(ctx, metricName)
			if !ok {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprintf("%d", value)))
		case "gauge":
			value, ok := m.Storage.GetGauge(ctx, metricName)
			if !ok {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprintf("%g", value)))
		default:
			{
				m.log.Debug("could not determine the type of metric")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		res.WriteHeader(http.StatusOK)
	}
	return getValueFn
}

func (m *MetricHandlerDB) UpdateNewMetrics(ctx context.Context) http.HandlerFunc {
	updateFn := func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			m.log.Debug("got request with bad method", zap.String("method", req.Method))
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
				m.log.Debug("cannot parse to int64", zap.Error(err))
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			err = m.Storage.AddCounterValue(ctx, metricName, valueint64)
			if err != nil {
				m.log.Debug("cannot get new counter value", zap.Error(err))
				res.WriteHeader(http.StatusNotFound)
				return
			}
		case "gauge":
			valuefloat64, err := strconv.ParseFloat(value, 64)
			if err != nil {
				m.log.Debug("cannot parse to float64", zap.Error(err))
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			err = m.Storage.AddGaugeValue(ctx, metricName, valuefloat64)
			if err != nil {
				m.log.Debug("cannot add new gauge value")
				res.WriteHeader(http.StatusNotFound)
				return
			}
		default:
			{
				m.log.Debug("could not determine the type of metric")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
	return updateFn
}

func (m *MetricHandlerDB) AddValueMetricsToJSON(ctx context.Context) http.HandlerFunc {
	addValueFn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			m.log.Debug("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// десериализуем запрос в структуру модели
		m.log.Debug("decoding request")
		var jsonMetrics models.Metrics
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&jsonMetrics); err != nil {
			m.log.Debug("cannot decode request JSON body", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		// заполняем модель ответа
		switch strings.ToLower(jsonMetrics.MType) {
		case "counter":
			value, ok := m.Storage.GetCounter(ctx, jsonMetrics.ID)
			if !ok {
				m.log.Debug("cannot get counter value")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			jsonMetrics.Delta = &value
		case "gauge":
			value, ok := m.Storage.GetGauge(ctx, jsonMetrics.ID)
			if !ok {
				m.log.Debug("cannot get gauge value")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			jsonMetrics.Value = &value
		default:
			{
				m.log.Debug("could not determine the type of metric")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		response, err := json.Marshal(jsonMetrics)
		if err != nil {
			m.log.Debug("cannot marshal to json", zap.Error(err))

		}
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(response)
	}
	return addValueFn
}
func (m *MetricHandlerDB) NewMetricsToJSON(ctx context.Context) http.HandlerFunc {
	newMetricsFn := func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			m.log.Debug("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// десериализуем запрос в структуру модели
		m.log.Debug("decoding request")
		var jsonMetrics models.Metrics
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&jsonMetrics); err != nil {
			m.log.Debug("cannot decode request JSON body", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		// заполняем модель ответа
		switch strings.ToLower(jsonMetrics.MType) {
		case "counter":
			err := m.Storage.AddCounterValue(ctx, jsonMetrics.ID, *jsonMetrics.Delta)
			if err != nil {
				m.log.Debug("cannot add new counter value")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			value, ok := m.Storage.GetCounter(ctx, jsonMetrics.ID)
			if !ok {
				m.log.Debug("cannot get counter value")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			jsonMetrics.Delta = &value
		case "gauge":
			err := m.Storage.AddGaugeValue(ctx, jsonMetrics.ID, *jsonMetrics.Value)
			if err != nil {
				m.log.Debug("cannot add new gauge value")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			value, ok := m.Storage.GetGauge(ctx, jsonMetrics.ID)
			if !ok {
				m.log.Debug("cannot get gauge value")
				res.WriteHeader(http.StatusNotFound)
				return
			}
			jsonMetrics.Value = &value
		default:
			{
				m.log.Debug("could not determine the type of metric")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		response, err := json.Marshal(jsonMetrics)
		if err != nil {
			m.log.Debug("cannot marshal to json", zap.Error(err))
			return
		}
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(response)
	}
	return newMetricsFn
}

func (m *MetricHandlerDB) PingDBConnect(ctx context.Context) http.HandlerFunc {
	pingFn := func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			m.log.Debug("got request with bad method", zap.String("method", req.Method))
			return
		}
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := m.pgDB.PingDB(ctx); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			res.WriteHeader(http.StatusOK)
		}
		_, err := res.Write([]byte("pong"))
		if err != nil {
			return
		}
	}
	return pingFn
}

func (m *MetricHandlerDB) UpdateNewMetricsBatch(ctx context.Context) http.HandlerFunc {
	updateMetricsFn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			m.log.Debug("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var metrics []models.Metrics
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&metrics); err != nil {
			m.log.Debug("cannot decode request JSON body", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		err := m.Storage.AddAllValue(ctx, metrics)
		if err != nil {
			m.log.Debug("cannot add new metrics value")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
	return updateMetricsFn
}
