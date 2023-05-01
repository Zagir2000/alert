package main

import (
	"log"
	"time"

	"github.com/Zagir2000/alert/internal/metricscollect"
	"github.com/go-resty/resty/v2"
)

type MyAPIError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	hostpath       = "http://localhost:8080/update/"
	reportInterval = 10
)

func sendMetrics(m *metricscollect.RuntimeMetrics) error {
	time.Sleep(reportInterval * time.Second)
	metrics := m.URLMetrics(hostpath)
	client := resty.New()
	var responseErr MyAPIError
	for _, url := range metrics {
		_, err := client.R().
			SetError(&responseErr).
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
