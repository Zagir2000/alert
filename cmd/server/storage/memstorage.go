package storage

import (
	"github.com/Zagir2000/alert/cmd/server/parser"
)

type gauge float64
type counter int64

type MemStorageUsage interface {
	CollectMetricsAndALerts(res string) error
}

type MemStorage struct {
	Gaugedata   map[string]gauge
	Counterdata map[string]counter
}

func (c *MemStorage) CollectMetricsAndALerts(res string) error {

	Metric, err := parser.Parseuri(res)
	if err != nil {
		return err
	}
	if Metric.Nametype == "counter" {
		c.Counterdata[Metric.Nametype] += counter(Metric.Valuecounter)
	} else {
		c.Gaugedata[Metric.Nametype] = gauge(Metric.Valuecounter)
	}
	return nil
}
