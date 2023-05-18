package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/metricscollect"
	"github.com/Zagir2000/alert/internal/server/handlers"
	"github.com/d5/tengo/assert"
)

func TestRunSendMetrics(t *testing.T) {
	type args struct {
		Metric metricscollect.RuntimeMetrics
		addr   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Pass test",
			args: args{
				Metric: metricscollect.RuntimeMetrics{
					RuntimeMemstats: map[string]float64{"Alloc": 123},
					PollCount:       5,
					RandomValue:     2,
				},
			},
			want: "context deadline exceeded",
		},
	}
	err := logger.Initialize("info")
	if err != nil {
		log.Println(err)
	}
	ts := httptest.NewServer(handlers.Router())
	u, err := url.Parse(ts.URL)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
	defer ts.Close()
	flag := NewFlagVarStruct()
	flag.parseFlags()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000*time.Millisecond))
			defer cancel()
			RunSendMetrics(test.args.Metric, ctx, cancel, u.Host)

			asd := fmt.Sprint(ctx.Err())
			assert.Equal(t, asd, test.want)

		})
	}
}
