package storage

import (
	"encoding/json"
	"os"
)

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

func (p *Producer) WriteMetrics(metricGauge *memStorage) error {
	err := p.encoder.Encode(&metricGauge)
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

func (c *Consumer) ReadMetrics() (*memStorage, error) {
	var metricsGaugeAndCounter memStorage
	if err := c.decoder.Decode(&metricsGaugeAndCounter); err != nil {
		return nil, err

	}
	return &metricsGaugeAndCounter, nil
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
	allMetrics := NewMemStorage()
	allMetrics.Gaugedata = m.GetAllGaugeValues()
	allMetrics.Counterdata = m.GetAllCounterValues()
	err = producer.WriteMetrics(allMetrics)
	if err != nil {
		return err
	}
	return nil
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
	metricsAll, err := consumer.ReadMetrics()
	if err != nil {
		return err
	}
	m.Counterdata = metricsAll.Counterdata
	m.Gaugedata = metricsAll.Gaugedata
	// сохраняем данные в файл

	return nil
}
