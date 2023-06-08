package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const query = `CREATE TABLE IF NOT EXISTS Metrics (co
	ID text not null,
	MTYPE text not null,
	DELTA bigint,
	VALUE double precision
	);`

type PostgresDB struct {
	db  *sql.DB
	log *zap.Logger
}

func (pgdb *PostgresDB) PingDB(ctx context.Context) error {
	err := pgdb.db.PingContext(ctx)
	return err
}
func InitDB(configDB string, log *zap.Logger) *PostgresDB {
	db, err := sql.Open("pgx", configDB)
	if err != nil {
		log.Error("Database initialization error")
		return nil

	}
	return &PostgresDB{db: db, log: log}
}

func (pgdb *PostgresDB) Close() {
	err := pgdb.db.Close()
	if err != nil {
		pgdb.log.Error("Error closing database connection: %v", zap.Error(err))
	}
}

func (session *PostgresDB) CreateTabel(ctx context.Context) error {
	_, err := session.db.ExecContext(ctx, query)
	if err != nil {
		session.log.Error("Error opening database connection: %v", zap.Error(err))
		return err
	}
	session.log.Info("Table metrics created successfully")
	return nil
}

func (m *PostgresDB) AddGaugeValue(name string, value float64) error {
	return nil
}

func (m *PostgresDB) AddCounterValue(name string, value int64) error {
	return nil
}

func (m *PostgresDB) GetGauge(name string) (float64, bool) {
	return 1, true
}

func (m *PostgresDB) GetCounter(name string) (int64, bool) {
	return 1, true
}

func (m *PostgresDB) GetAllGaugeValues() map[string]float64 {
	return nil
}

func (m *PostgresDB) GetAllCounterValues() map[string]int64 {
	return nil
}

func (m *PostgresDB) LoadMetricsJSON(metricsFile *memStorage) {
	return
}
