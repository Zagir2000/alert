package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Zagir2000/alert/internal/metricscollect"
	"github.com/go-resty/resty/v2"
)

type MyAPIError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

var flagRunAddr string
var reportInterval int
var pollInterval int

func parseFlags() {
	// как аргумент -a со значением :8080 по умолчанию
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")

	// частота отправки метрик на сервер
	flag.IntVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")

	//частота опроса метрик из пакета
	flag.IntVar(&pollInterval, "p", 2, "frequency of polling metrics from the package")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportIntervalInt, err := strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatalln(err)
		}
		reportInterval = envReportIntervalInt
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envPollIntervalInt, err := strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatalln(err)
		}
		pollInterval = envPollIntervalInt
	}
}

func sendMetrics(m *metricscollect.RuntimeMetrics) error {

	time.Sleep(time.Duration(reportInterval) * time.Millisecond)
	metrics := m.URLMetrics(flagRunAddr)
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
	parseFlags()
	Metric := metricscollect.PollIntervalPin(pollInterval)
	Metric.AddValueMetric()
	go Metric.NewСollect()
	for {
		err := sendMetrics(&Metric)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
