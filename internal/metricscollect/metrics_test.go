package metricscollect

import (
	"reflect"
	"testing"
	"time"

	"github.com/Zagir2000/alert/internal/models"
)

var value1 float64 = 1

var value2 float64 = 2

var value3 int64 = 3

var value4 float64 = 4

func TestRuntimeMetricsURLMetrics(t *testing.T) {
	type fields struct {
		RuntimeMemstats map[string]float64
		PollCount       int64
		RandomValue     float64
		pollInterval    time.Duration
	}
	type args struct {
		hostpath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.Metrics
	}{
		{
			name:   "Pass test",
			fields: fields{RuntimeMemstats: map[string]float64{"Test1": value1, "Test2": value2}, PollCount: value3, RandomValue: value4},
			args:   args{hostpath: "localhost:8080"},
			want: []models.Metrics{
				{
					ID:    "Test1",
					MType: "gauge",
					Value: &value1,
				},
				{
					ID:    "Test2",
					MType: "gauge",
					Value: &value2,
				},

				{
					ID:    "RandomValue",
					MType: "gauge",
					Value: &value4,
				},
				{
					ID:    "PollCount",
					MType: "counter",
					Delta: &value3,
				},
			},
		},
		{
			name:   "Failed test",
			fields: fields{RuntimeMemstats: map[string]float64{"asd_asd": 1, "asd_asd2": 2}, PollCount: 1231233, RandomValue: 1023123},
			args:   args{hostpath: "localhost:8080"},
			want:   []models.Metrics{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RuntimeMetrics{
				RuntimeMemstats: tt.fields.RuntimeMemstats,
				PollCount:       tt.fields.PollCount,
				RandomValue:     tt.fields.RandomValue,
				pollInterval:    tt.fields.pollInterval,
			}
			if got := m.URLMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RuntimeMetrics.URLMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
