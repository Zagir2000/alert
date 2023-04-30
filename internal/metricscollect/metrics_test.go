package metricscollect

import (
	"testing"
	"time"
)

func TestRuntimeMetrics_AddValueMetric(t *testing.T) {
	type fields struct {
		RuntimeMemstats map[string]float64
		PollCount       counter
		RandomValue     gauge
		pollInterval    time.Duration
	}
	tests := []struct {
		name   string
		fields fields
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
			m.AddValueMetric()
		})
	}
}
