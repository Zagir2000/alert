package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Zagir2000/alert/internal/agent/metricscollect"
	"golang.org/x/sync/errgroup"
)

func main() {

	flag := NewFlagVarStruct()
	err := flag.parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	g := new(errgroup.Group)
	Metric := metricscollect.IntervalPin(flag.pollInterval)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	// создаем буферизованный канал для отправки результатов
	jobs := make(chan []byte, 20)
	go Metric.NewСollectMetricGopsutil(ctx, cancel, jobs)
	go Metric.NewСollect(ctx, cancel, jobs)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(flag.pollInterval) * time.Second):
			for w := 1; w <= flag.rateLimit; w++ {
				fmt.Println(w)
				g.Go(func() error {
					var mx sync.Mutex
					mx.Lock()
					err := metricscollect.SendMetricsGor(jobs, flag.runAddr, flag.secretKey)
					if err != nil {
						return err
					}
					return nil
				})
			}
		}
		if err := g.Wait(); err != nil {
			cancel()
		}
	}

}
