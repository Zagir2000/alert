package storage

import (
	"github.com/Zagir2000/alert/internal/parser"
)

type Gauge float64
type Counter int64

type Repository interface {
	CollectMetricsAndALerts(res string) error
}

type MemStorage struct {
	Gaugedata   map[string]Gauge
	Counterdata map[string]Counter
}

func (c *MemStorage) CollectMetricsAndALerts(res string) error {

	Metric, err := parser.Parseuri(res)
	if err != nil {
		return err
	}
	if Metric.Nametype == "counter" {
		c.Counterdata[Metric.Nametype] += Counter(Metric.Valuecounter)
	} else {
		c.Gaugedata[Metric.Nametype] = Gauge(Metric.Valuecounter)
	}
	return nil
}
