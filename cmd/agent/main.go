package main

import (
	"log"
	"time"

	"github.com/Zagir2000/alert/internal/metricscollect"
	"github.com/go-resty/resty/v2"
)

const (
	hostpath       = "http://localhost:8080/update/"
	reportInterval = 10
)

func sendMetrics(m *metricscollect.RuntimeMetrics) error {
	time.Sleep(reportInterval * time.Second)
	metrics := m.URLMetrics(hostpath)
	for _, url := range metrics {
		client := resty.New()
		_, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(url)
		return err
	}
	return nil
}
func main() {
	Metric := metricscollect.PollIntervalPin()
	Metric.AddValueMetric()
	go Metric.NewСollect()
	for {
		err := sendMetrics(&Metric)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
