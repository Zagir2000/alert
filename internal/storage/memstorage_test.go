package storage

import (
	"context"
	"testing"
)

func TestAddGaugeValueStorage(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		value  float64
		want   *memStorage
	}{
		{
			name:   "Positive test",
			metric: "metricTest",
			value:  1.21,
			want: &memStorage{
				Gaugedata:   map[string]float64{"metricTest": 1.21},
				Counterdata: nil,
			},
		},
		{
			name:   "Failed test",
			metric: "metricTest",
			value:  78,
			want: &memStorage{
				Gaugedata:   map[string]float64{"metricTest": 1.21},
				Counterdata: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemStorage()
			got.AddGaugeValue(context.Background(), tt.metric, tt.value)
			if got == tt.want {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCounterStorage(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		value  int64
		want   *memStorage
	}{
		{
			name:   "Positive test",
			metric: "metricTest",
			value:  200,
			want: &memStorage{
				Counterdata: map[string]int64{"metricTest": 200},
				Gaugedata:   nil,
			},
		},
		{
			name:   "Failed test",
			metric: "metricTest",
			value:  78,
			want: &memStorage{
				Gaugedata: nil,
				Counterdata: map[string]int64{
					"metricTest": 123,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemStorage()
			got.AddCounterValue(context.Background(), tt.metric, tt.value)
			if got == tt.want {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
