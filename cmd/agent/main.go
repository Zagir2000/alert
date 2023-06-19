package main

import (
	"context"
	"log"

	"github.com/Zagir2000/alert/internal/agent/metricscollect"
)

func main() {
	flag := NewFlagVarStruct()
	err := flag.parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	Metric := metricscollect.IntervalPin(flag.pollInterval, flag.reportInterval)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go Metric.NewСollect(ctx, cancel)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := Metric.SendMetrics(flag.runAddr)
			if err != nil {
				log.Println("Error in send metrics:", err)
				cancel()
			}
		}
	}

}
