package storage

import (
	"context"
	"time"

	"github.com/Zagir2000/alert/internal/models"
	"go.uber.org/zap"
)

// Интерфейс который имплементирует все метода для сохранения в бд.
type Repository interface {
	AddGaugeValue(ctx context.Context, name string, value float64) error
	AddCounterValue(ctx context.Context, name string, value int64) error
	GetGauge(ctx context.Context, name string) (float64, bool)
	GetCounter(ctx context.Context, name string) (int64, bool)
	GetAllGaugeValues(ctx context.Context) map[string]float64
	GetAllCounterValues(ctx context.Context) map[string]int64
	AddAllValue(ctx context.Context, metrics []models.Metrics) error
}

// Инициализируем базу данных.
func NewStorage(ctx context.Context, migratePath string, log *zap.Logger, fileStoragePath string, restore bool, storeIntervall int, postgresDSN string) (Repository, *PostgresDB, error) {
	if postgresDSN != "" {
		DB, err := InitDB(postgresDSN, log, migratePath)
		if err != nil {
			log.Error("Error in initialization db", zap.Error(err))
			return nil, nil, err
		}
		return DB, DB, nil
	}

	if restore {
		memStorage := NewMemStorage()
		err := MetricsLoadJSON(fileStoragePath, memStorage)
		if err != nil {
			log.Error("failed to load file", zap.Error(err))
			return nil, nil, err
		}
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		go func() {
			select {
			case <-ctx.Done():
				return
			default:
				for {
					err = MetricsSaveJSON(fileStoragePath, memStorage)
					if err != nil {
						log.Error("failed to save file", zap.Error(err))
						cancel()
					}
					time.Sleep(time.Duration(storeIntervall) * time.Second)
				}
			}
		}()
		return memStorage, nil, nil
	}
	return NewMemStorage(), nil, nil
}
