package storage

import (
	"testing"
)

func TestSetGaugeStorage(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		value  float64
		want   *MemStorage
	}{
		{
			name:   "Positive test",
			metric: "metricTest",
			value:  1.21,
			want: &MemStorage{
				Gaugedata:   map[string]float64{"metricTest": 1.21},
				Counterdata: nil,
			},
		},
		{
			name:   "Failed test",
			metric: "metricTest",
			value:  78,
			want: &MemStorage{
				Gaugedata:   map[string]float64{"metricTest": 1.21},
				Counterdata: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemStorage()
			got.SetGauge(tt.metric, tt.value)
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
		want   *MemStorage
	}{
		{
			name:   "Positive test",
			metric: "metricTest",
			value:  200,
			want: &MemStorage{
				Counterdata: map[string]int64{"metricTest": 200},
				Gaugedata:   nil,
			},
		},
		{
			name:   "Failed test",
			metric: "metricTest",
			value:  78,
			want: &MemStorage{
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
			got.SetCounter(tt.metric, tt.value)
			if got == tt.want {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
