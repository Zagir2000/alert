package storage

import (
	"errors"

	"go.uber.org/zap"
)

type memStorage struct {
	Gaugedata   map[string]float64
	Counterdata map[string]int64
	log         *zap.Logger
}

func NewMemStorage() *memStorage {
	return &memStorage{
		Gaugedata:   make(map[string]float64),
		Counterdata: make(map[string]int64),
	}
}

func (m *memStorage) AddGaugeValue(name string, value float64) error {
	m.Gaugedata[name] = value
	valuenew, ok := m.Gaugedata[name]
	if !ok && value == valuenew {
		return errors.New("failed to add gauge value")
	}
	return nil
}

func (m *memStorage) AddCounterValue(name string, value int64) error {
	if value < 0 {
		return errors.New("counter cannot decrease in value")
	}
	m.Counterdata[name] += value
	if m.Counterdata[name] == 0 && value != 0 {
		return errors.New("counter is overflow")
	}
	_, ok := m.Counterdata[name]
	if !ok {
		return errors.New("failed to add counter value")
	}
	return nil
}

func (m *memStorage) GetGauge(name string) (float64, bool) {
	value, ok := m.Gaugedata[name]
	return value, ok
}

func (m *memStorage) GetCounter(name string) (int64, bool) {
	value, ok := m.Counterdata[name]
	return value, ok
}

func (m *memStorage) GetAllGaugeValues() map[string]float64 {
	return m.Gaugedata
}

func (m *memStorage) GetAllCounterValues() map[string]int64 {
	return m.Counterdata
}

func (m *memStorage) LoadMetricsJSON(metricsFile *memStorage) {
	m.Counterdata = metricsFile.Counterdata
	m.Gaugedata = metricsFile.Gaugedata
}
