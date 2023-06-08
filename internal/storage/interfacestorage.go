package storage

import (
	"go.uber.org/zap"
)

type Repository interface {
	AddGaugeValue(name string, value float64) error
	AddCounterValue(name string, value int64) error
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGaugeValues() map[string]float64
	GetAllCounterValues() map[string]int64
	LoadMetricsJSON(metricGaugeFile *memStorage)
}

func MetricHandler(fileStoragePath string, restore bool, storeInterval int, log *zap.Logger, postgresDSN string) (Repository, *PostgresDB) {
	if postgresDSN != "" {
		Db := InitDB(postgresDSN, log)
		return Db, Db
	}

	if restore {
		memStorage := NewMemStorage()
		err := memStorage.MetricsLoadJSON(fileStoragePath)
		if err != nil {
			log.Error("failed to load file", zap.Error(err))
		}
		return memStorage, nil
	}
	return NewMemStorage(), nil
}
