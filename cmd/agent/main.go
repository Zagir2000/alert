package main

import (
	"context"
	"fmt"
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
	// создаем буферизованный канал для отправки результатов
	jobs := make(chan []byte, 1)

	go Metric.NewСollect(ctx, cancel, jobs)
	go Metric.NewСollectMetricGopsutil(ctx, cancel, jobs)
	for {

		for w := 1; w <= flag.rateLimit; w++ {
			time.Sleep(time.Duration(flag.pollInterval) * time.Second)
			fmt.Println(w)
			go Metric.SendMetricsGor(ctx, cancel, jobs, flag.runAddr, flag.secretKey)

		}

	}

}
