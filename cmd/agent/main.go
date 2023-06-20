package main

import (
	"context"
	"log"
	"time"

	"github.com/Zagir2000/alert/internal/agent/metricscollect"
)

func main() {
	flag := NewFlagVarStruct()
	err := flag.parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	Metric := metricscollect.IntervalPin(flag.pollInterval)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go Metric.NewСollect(ctx, cancel)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(flag.reportInterval) * time.Second)
			err := Metric.SendMetrics(flag.runAddr, flag.secretKey)
			if err != nil {
				log.Println("Error in send metrics:", err)
				cancel()
			}
		}
	}

}
