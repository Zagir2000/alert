package storage

import (
	"github.com/Zagir2000/alert/internal/parser"
)

type Repository interface {
	CollectMetricsAndALerts(res string) error
}

type MemStorage struct {
	Gaugedata   map[string]float64
	Counterdata map[string]int64
}

func (c *MemStorage) CollectMetricsAndALerts(res string) error {

	Metric, err := parser.Parseuri(res)
	if err != nil {
		return err
	}
	if Metric.Type == "counter" {
		c.Counterdata[Metric.Nametype] = Metric.Valuecounter
	} else {
		c.Gaugedata[Metric.Nametype] = Metric.Valuegauge
	}
	return nil
}
