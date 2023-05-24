package storage

import (
	"encoding/json"
	"os"
)

type metricsFile struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
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

func (p *Producer) WriteMetrics(metric *metricsFile) error {
	return p.encoder.Encode(&metric)
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Consumer{file: file, decoder: json.NewDecoder(file)}, nil
}

func (c *Consumer) ReadMetrics() (*metricsFile, error) {
	var metrics metricsFile
	if err := c.decoder.Decode(&metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func MetricsSaveJson(fname string, m *memStorage) error {
	producer, err := NewProducer(fname)
	if err != nil {
		return err
	}
	defer producer.Close()
	// сохраняем данные в файл
	metrics := &metricsFile{}
	metrics.Counter = m.GetAllCounterValues()
	metrics.Gauge = m.GetAllGaugeValues()

	return producer.WriteMetrics(metrics)
}

func MetricsLoadJSON(fname string, m *memStorage) error {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return nil
	}
	consumer, err := NewConsumer(fname)
	if err != nil {
		return err
	}
	defer consumer.Close()
	metricsFile, err := consumer.ReadMetrics()
	if err != nil {
		return err
	}
	m.LoadMetricsJSON(metricsFile)
	// сохраняем данные в файл

	return nil
}
