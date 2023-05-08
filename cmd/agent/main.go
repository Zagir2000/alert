package main

import (
	"log"

	"github.com/Zagir2000/alert/internal/metricscollect"
)

func main() {
	runAddr, pollInterval, reportInterval := parseFlags()
	Metric := metricscollect.IntervalPin(pollInterval, reportInterval)
	go Metric.NewСollect()

	for {
		err := Metric.SendMetrics(runAddr)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
