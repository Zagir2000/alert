package metricscollect

import (
	"reflect"
	"testing"
	"time"
)

func TestRuntimeMetrics_URLMetrics(t *testing.T) {
	type fields struct {
		RuntimeMemstats map[string]float64
		PollCount       counter
		RandomValue     gauge
		pollInterval    time.Duration
	}
	type args struct {
		hostpath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "Pass test",
			fields: fields{RuntimeMemstats: map[string]float64{"Test1": 1, "Test2": 2}, PollCount: 3, RandomValue: 4},
			args:   args{hostpath: "localhost:8080"},
			want: []string{"http://localhost:8080/update/gauge/Test1/1.000000",
				"http://localhost:8080/update/gauge/Test2/2.000000",
				"http://localhost:8080/update/gauge/RandomValue/4.000000",
				"http://localhost:8080/update/counter/PollCount/3"},
		},
		{
			name:   "Failed test",
			fields: fields{RuntimeMemstats: map[string]float64{"asd_asd": 1, "asd_asd2": 2}, PollCount: 1231233, RandomValue: 1023123},
			args:   args{hostpath: "localhost:8080"},
			want: []string{"http://localhost:8080/update/gauge/asd_asd/1.000000",
				"http://localhost:8080/update/gauge/asd_asd2/2.000000",
				"http://localhost:8080/update/gauge/RandomValue/1023123.000000",
				"http://localhost:8080/update/counter/PollCount/1231233"},
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
			if got := m.URLMetrics(tt.args.hostpath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RuntimeMetrics.URLMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
