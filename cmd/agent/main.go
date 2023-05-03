package main

import (
	"log"

	"github.com/Zagir2000/alert/internal/metricscollect"
)

func main() {
	parseFlags()
	Metric := metricscollect.PollIntervalPin(pollInterval)
	Metric.AddValueMetric()
	go Metric.NewСollect()
	for {
		err := Metric.SendMetrics(flagRunAddr)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
