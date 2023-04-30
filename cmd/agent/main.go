package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/Zagir2000/alert/internal/metricscollect"
)

const (
	hostpath       = "http://localhost:8080/update/"
	reportInterval = 10
)

func sendMetrics(m *metricscollect.RuntimeMetrics) error {
	time.Sleep(reportInterval * time.Second)
	metrics := m.UrlMetrics(hostpath)
	for _, url := range metrics {
		req, err := http.Post(url, "text/plain", bytes.NewBuffer([]byte{}))
		if err != nil {
			return err
		}
		defer req.Body.Close()
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
