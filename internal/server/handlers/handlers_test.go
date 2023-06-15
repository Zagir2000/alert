package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zagir2000/alert/internal/server/logger"
	"github.com/Zagir2000/alert/internal/server/storage"
	"github.com/d5/tengo/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricHandler_MainPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Pass test",
			url:  "/update/counter/someMetric/1",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Failed test 1",
			url:  "/update/counter/someMetric/b",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Failed test 2",
			url:  "/update/counter/1",
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Failed test 3",
			url:  "/update/counter/",
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
	memStorage := storage.NewMemStorage()
	newHandStruct := MetricHandlerNew(memStorage, logger, nil)
	r := Router(context.Background(), newHandStruct)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("POST", tt.url, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			assert.Equal(t, w.Code, tt.want.code)
			assert.Equal(t, w.Code, tt.want.code)
		})
	}
}
