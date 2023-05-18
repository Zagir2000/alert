package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/server/handlers"
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
				contentType: "",
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
	flagStruct := NewFlagVarStruct()
	flagStruct.parseFlags()
	err := logger.Initialize(flagStruct.logLevel)
	if err != nil {
		log.Println(err)
	}
	ts := httptest.NewServer(handlers.Router())
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
