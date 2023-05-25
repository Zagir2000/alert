package storage

import (
	"encoding/json"
	"os"
)

type metricsFileGauge struct {
	Gauge map[string]float64 `json:"gauge"`
}
type metricsFileCounter struct {
	Counter map[string]int64 `json:"counter"`
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

type Consumer struct {
	file    *os.File // файл для чтения
	decoder *json.Decoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteMetrics(metricGauge *metricsFileGauge, metricCounter *metricsFileCounter) error {
	err := p.encoder.Encode(&metricGauge)
	if err != nil {
		return err
	}
	err = p.encoder.Encode(&metricCounter)
	if err != nil {
		return err
	}
	return nil
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Consumer{file: file, decoder: json.NewDecoder(file)}, nil
}

func (c *Consumer) ReadMetrics() (*metricsFileGauge, *metricsFileCounter, error) {
	var metricsGauge metricsFileGauge
	var metricsCounter metricsFileCounter
	if err := c.decoder.Decode(&metricsGauge); err != nil {
		return nil, nil, err
	}
	if err := c.decoder.Decode(&metricsCounter); err != nil {
		return nil, nil, err
	}
	return &metricsGauge, &metricsCounter, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func MetricsSaveJSON(fname string, m *memStorage) error {
	producer, err := NewProducer(fname)
	if err != nil {
		return err
	}
	defer producer.Close()
	// сохраняем данные в файл
	metricsGauge := &metricsFileGauge{}
	metricsCounter := &metricsFileCounter{}
	metricsCounter.Counter = m.GetAllCounterValues()
	metricsGauge.Gauge = m.GetAllGaugeValues()
	err = producer.WriteMetrics(metricsGauge, metricsCounter)
	if err != nil {
		return err
	}
	return err
}

func (m *memStorage) MetricsLoadJSON(fname string) error {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return nil
	}
	consumer, err := NewConsumer(fname)
	if err != nil {
		return err
	}
	defer consumer.Close()
	metricsGauge, MetricsCounter, err := consumer.ReadMetrics()
	if err != nil {
		return err
	}
	m.LoadMetricsJSON(metricsGauge, MetricsCounter)
	// сохраняем данные в файл

	return nil
}
