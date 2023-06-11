package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zagir2000/alert/internal/handlers"
	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/storage"
	"github.com/d5/tengo/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "positive test #1",
			args: "/update/counter/metric/1",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test: wrong value #1",
			args: "/update/counter/metric/b",
			want: want{
				code:        400,
				contentType: "application/x-gzip",
			},
		},
		{
			name: "negative test: missing metric name #2",
			args: "/update/counter/",
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	logger, err := logger.InitializeLogger("info")
	if err != nil {
		log.Println(err)
	}
	m := storage.NewMemStorage()
	newHandStruct := handlers.MetricHandlerNew(m, logger, nil)
	ts := httptest.NewServer(handlers.Router(context.Background(), newHandStruct))
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL+test.args, nil)
			require.NoError(t, err)
			resp, err := ts.Client().Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, test.want.code)
			assert.Equal(t, resp.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
