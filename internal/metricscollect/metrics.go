package metricscollect

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	counterMetric   string = "counter"
	gaugeMetric     string = "gauge"
	RandomValueName string = "RandomValue"
	PollCountName   string = "PollCount"
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
	return RuntimeMetrics{pollInterval: time.Duration(pollIntervalFlag), reportInterval: time.Duration(reportIntervalFlag)}
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
	m.RandomValue = rand.Float64()
	if m.PollCount < 0 {
		return errors.New("counter is negative number")
	}
	m.PollCount += 1
	m.RuntimeMemstats = mapstats
	time.Sleep(m.pollInterval * time.Second)
	return nil
}

func (m *RuntimeMetrics) URLMetrics(hostpath string) []string {
	urls := make([]string, 0, len(m.RuntimeMemstats)+2)
	for i, k := range m.RuntimeMemstats {
		s := fmt.Sprintf("%f", k)
		URL := strings.Join([]string{"http:/", hostpath, "update", gaugeMetric, i, s}, "/")
		urls = append(urls, URL)
	}
	s := fmt.Sprintf("%f", m.RandomValue)

	URLRandomGuage := strings.Join([]string{"http:/", hostpath, "update", gaugeMetric, RandomValueName, s}, "/")
	c := fmt.Sprintf("%d", m.PollCount)

	URLCount := strings.Join([]string{"http:/", hostpath, "update", counterMetric, PollCountName, c}, "/")
	urls = append(urls, URLRandomGuage, URLCount)
	return urls
}

func (m *RuntimeMetrics) SendMetrics(hostpath string) error {

	time.Sleep(m.reportInterval * time.Second)
	metrics := m.URLMetrics(hostpath)
	client := resty.New()
	var responseErr SendMetricsError
	for _, url := range metrics {
		_, err := client.R().
			SetError(&responseErr).
			SetHeader("Content-Type", "text/plain").
			Post(url)

		if err != nil {
			return err
		}
	}
	return nil
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
