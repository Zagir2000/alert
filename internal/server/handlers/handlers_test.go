package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zagir2000/alert/internal/logger"
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
	err := logger.Initialize("info")
	if err != nil {
		log.Println(err)
	}
	r := Router()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("POST", tt.url, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			assert.Equal(t, w.Code, tt.want.code)
			s := w.Header()
			fmt.Println(s.Get("Content-Type"))
			assert.Equal(t, s.Get("Content-Type"), tt.want.contentType)
		})
	}
}
