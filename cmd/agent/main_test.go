package main

import (
	"testing"

	"github.com/Zagir2000/alert/internal/metricscollect"
)

func Test_sendMetrics(t *testing.T) {
	type args struct {
		m *metricscollect.RuntimeMetrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test send request", // описываем каждый тест:
			args: args{&metricscollect.RuntimeMetrics{RuntimeMemstats: map[string]float64{"Alloc": 1,
				"BuckHashSys": 2}, PollCount: 1, RandomValue: 1}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendMetrics(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("sendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}