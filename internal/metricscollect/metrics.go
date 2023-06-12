package metricscollect

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Zagir2000/alert/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/johncgriffin/overflow"
)

const (
	counterMetric   string = "counter"
	gaugeMetric     string = "gauge"
	randomValueName string = "RandomValue"
	pollCountName   string = "PollCount"
	contentType     string = "application/json"
	compressType    string = "gzip"
)

type RuntimeMetrics struct {
	RuntimeMemstats map[string]float64
	PollCount       int64
	RandomValue     float64
	pollInterval    time.Duration
	reportInterval  time.Duration
}

type SendMetricsError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func IntervalPin(pollIntervalFlag int, reportIntervalFlag int) RuntimeMetrics {
	return RuntimeMetrics{pollInterval: time.Duration(pollIntervalFlag), reportInterval: time.Duration(reportIntervalFlag), RuntimeMemstats: make(map[string]float64), PollCount: 0, RandomValue: 0}
}

func (m *RuntimeMetrics) AddValueMetric() error {

	mapstats := make(map[string]float64)
	var RtMetrics runtime.MemStats
	runtime.ReadMemStats(&RtMetrics)
	mapstats["Alloc"] = float64(RtMetrics.Alloc)
	mapstats["BuckHashSys"] = float64(RtMetrics.BuckHashSys)
	mapstats["Frees"] = float64(RtMetrics.Frees)
	mapstats["GCCPUFraction"] = float64(RtMetrics.GCCPUFraction)
	mapstats["GCSys"] = float64(RtMetrics.GCSys)
	mapstats["HeapAlloc"] = float64(RtMetrics.HeapAlloc)
	mapstats["HeapIdle"] = float64(RtMetrics.HeapIdle)
	mapstats["HeapInuse"] = float64(RtMetrics.HeapInuse)
	mapstats["HeapObjects"] = float64(RtMetrics.HeapObjects)
	mapstats["HeapReleased"] = float64(RtMetrics.HeapReleased)
	mapstats["HeapSys"] = float64(RtMetrics.HeapSys)
	mapstats["LastGC"] = float64(RtMetrics.LastGC)
	mapstats["Lookups"] = float64(RtMetrics.Lookups)
	mapstats["MCacheInuse"] = float64(RtMetrics.MCacheInuse)
	mapstats["MCacheSys"] = float64(RtMetrics.MCacheSys)
	mapstats["MSpanInuse"] = float64(RtMetrics.MSpanInuse)
	mapstats["MSpanSys"] = float64(RtMetrics.MSpanSys)
	mapstats["Mallocs"] = float64(RtMetrics.Mallocs)
	mapstats["NextGC"] = float64(RtMetrics.NextGC)
	mapstats["NumForcedGC"] = float64(RtMetrics.NumForcedGC)
	mapstats["NumGC"] = float64(RtMetrics.NumGC)
	mapstats["OtherSys"] = float64(RtMetrics.OtherSys)
	mapstats["PauseTotalNs"] = float64(RtMetrics.PauseTotalNs)
	mapstats["StackInuse"] = float64(RtMetrics.StackInuse)
	mapstats["StackSys"] = float64(RtMetrics.StackSys)
	mapstats["Sys"] = float64(RtMetrics.Sys)
	mapstats["TotalAlloc"] = float64(RtMetrics.TotalAlloc)
	m.RandomValue = rand.Float64() * 10000
	if m.PollCount < 0 {
		return errors.New("counter is negative number")
	}

	checkCounterInOverflow, ok := overflow.Add64(m.PollCount, 1)

	if !ok {
		m.PollCount = 0
		return errors.New("counter is overflow")
	}
	m.PollCount = checkCounterInOverflow
	m.RuntimeMemstats = mapstats
	time.Sleep(m.pollInterval * time.Second)
	return nil
}
func (m *RuntimeMetrics) jsonMetricsToBatch() []byte {
	var metrics []models.Metrics
	for k, v := range m.RuntimeMemstats {
		jsonGauge := &models.Metrics{
			ID:    k,
			MType: gaugeMetric,
			Value: &v,
		}
		metrics = append(metrics, *jsonGauge)
	}
	URLRandomGauge := &models.Metrics{
		ID:    randomValueName,
		MType: gaugeMetric,
		Value: &m.RandomValue,
	}
	metrics = append(metrics, *URLRandomGauge)
	URLCount := &models.Metrics{
		ID:    pollCountName,
		MType: counterMetric,
		Delta: &m.PollCount,
	}
	metrics = append(metrics, *URLCount)
	out, err := json.Marshal(metrics)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func (m *RuntimeMetrics) SendMetrics(hostpath string) error {

	time.Sleep(m.reportInterval * time.Second)
	client := resty.New()
	var responseErr SendMetricsError
	url := strings.Join([]string{"http:/", hostpath, "updates/"}, "/")
	out := m.jsonMetricsToBatch()
	res, err := gzipCompress(out)
	if err != nil {
		return err
	}
	_, err = client.R().
		SetError(&responseErr).
		SetHeader("Content-Encoding", compressType).
		SetHeader("Content-Type", contentType).
		SetBody(res).
		Post(url)
	if err != nil {
		if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, io.EOF) {
			for _, k := range models.TimeConnect {
				time.Sleep(k)
				fmt.Println(k)
				_, err := client.R().
					SetHeader("Content-Type", contentType).
					SetHeader("Content-Encoding", compressType).
					SetBody(res).
					Post(url)
				if err == nil {
					break
				}
			}
		}
	}
	return err
}

func (m *RuntimeMetrics) NewСollect(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := m.AddValueMetric()
			if err != nil {
				log.Println("Error in collect metrics:", err)
			}
		}

	}
}

func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}
