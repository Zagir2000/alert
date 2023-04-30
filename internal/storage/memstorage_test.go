package storage

import "testing"

func TestMemStorage_CollectMetricsAndALerts(t *testing.T) {
	type fields struct {
		Gaugedata   map[string]float64
		Counterdata map[string]int64
	}
	type args struct {
		res string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test collect counter method", // описываем каждый тест:
			fields:  fields{Counterdata: map[string]int64{"NameMetric": 0}},
			args:    args{res: "/update/counter/Namemetric/0"},
			wantErr: false,
		},
		{
			name:    "Test collect gauge method", // описываем каждый тест:
			fields:  fields{Gaugedata: map[string]float64{"NameMetric": 0}},
			args:    args{res: "/update/gauge/Namemetric/0"},
			wantErr: false,
		},
		{
			name:    "Err collect gauge method", // описываем каждый тест:
			fields:  fields{Gaugedata: map[string]float64{"NameMetric": 0}},
			args:    args{res: "/update/none/Namemetric/0"},
			wantErr: true,
		},
		{
			name:    "Err collect counter method", // описываем каждый тест:
			fields:  fields{Counterdata: map[string]int64{"NameMetric": 0}},
			args:    args{res: "/update/none/Namemetric/0"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemStorage{
				Gaugedata:   tt.fields.Gaugedata,
				Counterdata: tt.fields.Counterdata,
			}
			if err := c.CollectMetricsAndALerts(tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.CollectMetricsAndALerts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
