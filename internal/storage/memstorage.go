package storage

import (
	"errors"
)

type Repository interface {
	AddGaugeValue(name string, value float64)
	AddCounterValue(name string, value int64) error
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGaugeValues() map[string]float64
	GetAllCounterValues() map[string]int64
}

type memStorage struct {
	Gaugedata   map[string]float64
	Counterdata map[string]int64
}

func NewMemStorage() *memStorage {
	return &memStorage{
		Gaugedata:   make(map[string]float64),
		Counterdata: make(map[string]int64),
	}
}

func (m *memStorage) AddGaugeValue(name string, value float64) {
	m.Gaugedata[name] = value
}

func (m *memStorage) AddCounterValue(name string, value int64) error {
	if value < 0 {
		return errors.New("counter cannot decrease in value")
	}
	m.Counterdata[name] += value
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
