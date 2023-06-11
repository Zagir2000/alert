package storage

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Repository interface {
	AddGaugeValue(name string, value float64) error
	AddCounterValue(name string, value int64) error
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGaugeValues() map[string]float64
	GetAllCounterValues() map[string]int64
}

func NewStorage(log *zap.Logger, fileStoragePath string, restore bool, storeIntervall int, postgresDSN string) (Repository, *PostgresDB) {
	if postgresDSN != "" {
		DB := InitDB(postgresDSN, log)
		DB.CreateTabel(context.Background())

		return DB, DB
	}

	if restore {
		memStorage := NewMemStorage()
		err := MetricsLoadJSON(fileStoragePath, memStorage)
		if err != nil {
			log.Error("failed to load file", zap.Error(err))
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
		return memStorage, nil
	}
	return NewMemStorage(), nil
}
