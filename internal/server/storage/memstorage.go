package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/Zagir2000/alert/internal/models"
)

// Структура для сохранения метрик в RAM.
type memStorage struct {
	Gaugedata   map[string]float64
	Counterdata map[string]int64
	rw          sync.RWMutex
}

// Инициализация структуры memStorage.
func NewMemStorage() *memStorage {
	return &memStorage{
		Gaugedata:   make(map[string]float64),
		Counterdata: make(map[string]int64),
	}
}

// Добавляем метрику gauge.
func (m *memStorage) AddGaugeValue(ctx context.Context, name string, value float64) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.Gaugedata[name] = value
	valuenew, ok := m.Gaugedata[name]
	if !ok && value == valuenew {
		return errors.New("failed to add gauge value")
	}
	return nil
}

// Добавляем метрику counter.
func (m *memStorage) AddCounterValue(ctx context.Context, name string, value int64) error {
	m.rw.Lock()
	defer m.rw.Unlock()
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

// Получаем метрику gauge.
func (m *memStorage) GetGauge(ctx context.Context, name string) (float64, bool) {
	value, ok := m.Gaugedata[name]
	return value, ok
}

// Получаем метрику counter.
func (m *memStorage) GetCounter(ctx context.Context, name string) (int64, bool) {
	value, ok := m.Counterdata[name]
	return value, ok
}

// Получаем все метрики gauge.
func (m *memStorage) GetAllGaugeValues(ctx context.Context) map[string]float64 {
	return m.Gaugedata
}

// Получаем все метрики counter.
func (m *memStorage) GetAllCounterValues(ctx context.Context) map[string]int64 {
	return m.Counterdata
}

// Получаем все метрики.
func (m *memStorage) AddAllValue(ctx context.Context, metrics []models.Metrics) error {
	for _, v := range metrics {
		// все изменения записываются в транзакцию
		if v.MType == "gauge" {
			err := m.AddGaugeValue(ctx, v.ID, *v.Value)
			if err != nil {
				return err
			}
		} else {
			err := m.AddCounterValue(ctx, v.ID, *v.Delta)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
