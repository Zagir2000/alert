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
	"go.uber.org/zap"
)

func TestMetricHandler_MainPage(t *testing.T) {
	type want struct {
		code        int
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
	newHandStruct := MetricHandlerNew(memStorage, nil)
	r := Router(context.Background(), logger, newHandStruct, "")
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

func Example() {
	//инициализируем логер
	log, _ := logger.InitializeLogger("info")

	ctx := context.Background()
	//папка с миграциями
	migrationsDir := "migrations"
	//значение, указывающее, следует ли загружать ранее сохраненные значения из указанного файла при запуске сервера
	restore := true
	//time interval according to which the current server servers are kept on disk
	//интервал времени, в течение которого текущие серверы сервера хранятся на диске
	storeIntervall := 300
	databaseDsn := "postgres://postgres:123456@localhost/metrics?sslmode=disable"
	memStorageInterface, postgresDB, err := storage.NewStorage(ctx, migrationsDir, log, migrationsDir, restore, storeIntervall, databaseDsn)
	if err != nil {
		log.Fatal("Error in create storage", zap.Error(err))
	}
	if postgresDB != nil {
		defer postgresDB.Close()
	}
	newHandStruct := MetricHandlerNew(memStorageInterface, postgresDB)
	// Выполняем операцию получения метрик.
	newHandStruct.GetAllMetrics(ctx, log)
	// Выполняем операцию добавления метрики, который использует json формат.
	newHandStruct.AddValueMetricsToJSON(ctx, log)
	// Выполняем операцию получения метрик.
	newHandStruct.GetNowValueMetrics(ctx, log)
	// Выполняем операцию добавления метрик, который использует json формат.
	newHandStruct.NewMetricsToJSON(ctx, log)
	// Выполняем операцию проверки соединения базы данных.
	newHandStruct.PingDBConnect(ctx, log)
	// Выполняем операцию обновления метрики.
	newHandStruct.UpdateNewMetrics(ctx, log)
	//Ключ для подписи хэша
	key := "secret"
	// Выполняем операцию обновления метрик.
	newHandStruct.UpdateNewMetricsBatch(ctx, log, key)
}
