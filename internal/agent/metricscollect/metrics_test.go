package metricscollect

import (
	"testing"
	"time"
)

func TestRuntimeMetrics_SendMetrics(t *testing.T) {
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
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RuntimeMetrics{
				RuntimeMemstats: tt.fields.RuntimeMemstats,
				PollCount:       tt.fields.PollCount,
				RandomValue:     tt.fields.RandomValue,
				pollInterval:    tt.fields.pollInterval,
			}
			if err := m.SendMetrics(tt.args.hostpath, ""); (err != nil) != tt.wantErr {
				t.Errorf("RuntimeMetrics.SendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
