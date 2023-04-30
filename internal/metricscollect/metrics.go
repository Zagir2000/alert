package metricscollect

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type gauge float64
type counter int64

const (
	counterMetric   string = "counter"
	gaugeMetric     string = "gauge"
	RandomValueName string = "RandomValue"
	PollCountName   string = "PollCount"
)

type RuntimeMetrics struct {
	RuntimeMemstats map[string]float64
	PollCount       counter
	RandomValue     gauge
	pollInterval    time.Duration
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

func PollIntervalPin() RuntimeMetrics {
	return RuntimeMetrics{pollInterval: 2}
}
func (m *RuntimeMetrics) AddValueMetric() {
	mapstats := make(map[string]float64)
	var RtMetrics runtime.MemStats
	runtime.ReadMemStats(&RtMetrics)
	for _, k := range gaugeMetrics {
		mapstats[k] = float64(RtMetrics.Alloc)
	}
	m.RandomValue = gauge(rand.Float64())
	m.PollCount += 1
	m.RuntimeMemstats = mapstats
	time.Sleep(m.pollInterval * time.Second)
}

func (m *RuntimeMetrics) UrlMetrics(hostpath string) []string {
	var urls []string
	for i, k := range m.RuntimeMemstats {
		s := fmt.Sprintf("%f", k)
		Url := hostpath + gaugeMetric + "/" + i + "/" + s
		urls = append(urls, Url)
	}
	s := fmt.Sprintf("%f", m.RandomValue)
	UrlRandomGuage := hostpath + gaugeMetric + "/" + RandomValueName + "/" + s
	c := fmt.Sprintf("%d", m.PollCount)
	UrlCount := hostpath + gaugeMetric + "/" + PollCountName + "/" + c
	urls = append(urls, UrlRandomGuage, UrlCount)
	return urls
}

func (m *RuntimeMetrics) NewСollect() {
	for {
		m.AddValueMetric()
	}
}
