package models

import "time"

// Структура для одной метрики
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// Время подключение к серверу
var TimeConnect []time.Duration = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
