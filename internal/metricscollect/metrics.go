package metricscollect

import (
	"fmt"
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

var gaugeMetrics = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc"}

func IntervalPin(pollIntervalFlag int, reportIntervalFlag int) RuntimeMetrics {
	return RuntimeMetrics{pollInterval: time.Duration(pollIntervalFlag), reportInterval: time.Duration(reportIntervalFlag)}
}
func (m *RuntimeMetrics) AddValueMetric() {
	mapstats := make(map[string]float64)
	var RtMetrics runtime.MemStats
	runtime.ReadMemStats(&RtMetrics)
	for _, k := range gaugeMetrics {
		mapstats[k] = float64(RtMetrics.Alloc)
	}
	m.RandomValue = rand.Float64()
	m.PollCount += 1
	m.RuntimeMemstats = mapstats
	time.Sleep(m.pollInterval * time.Second)
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

func (m *RuntimeMetrics) NewСollect() {
	for {
		m.AddValueMetric()

	}
}
