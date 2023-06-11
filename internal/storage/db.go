package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const query = `CREATE TABLE IF NOT EXISTS Metrics (
	ID TEXT,
	MTYPE TEXT,
	DELTA INTEGER,
	VALUE DOUBLE PRECISION
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

func (pgdb *PostgresDB) CreateTabel(ctx context.Context) error {
	_, err := pgdb.db.ExecContext(ctx, query)
	if err != nil {
		pgdb.log.Error("Error opening database connection: %v", zap.Error(err))
		return err
	}
	pgdb.log.Info("Table metrics created successfully")
	return nil
}

func (pgdb *PostgresDB) AddGaugeValue(name string, value float64) error {
	_, err := pgdb.db.ExecContext(context.Background(),
		`INSERT INTO metrics (ID,MTYPE,VALUE) VALUES ($1, $2, $3);`, name, "gauge", value)
	return err
}

func (pgdb *PostgresDB) AddCounterValue(name string, value int64) error {
	_, err := pgdb.db.ExecContext(context.Background(),
		`INSERT INTO metrics (ID,MTYPE,DELTA) VALUES ($1, $2, $3);`, name, "counter", value)
	return err
}

func (pgdb *PostgresDB) GetGauge(name string) (float64, bool) {
	var value float64
	row := pgdb.db.QueryRowContext(context.Background(),
		"SELECT metrics.VALUE  FROM metrics WHERE metrics.ID=$1", name)

	err := row.Scan(&value)
	if err != nil {
		pgdb.log.Error("Error in get gauage value", zap.Error(err))
		return 0, false
	}
	return value, true
}

func (pgdb *PostgresDB) GetCounter(name string) (int64, bool) {
	var value int64
	row := pgdb.db.QueryRowContext(context.Background(),
		"SELECT metrics.DELTA  FROM metrics WHERE metrics.ID=$1", name)
	err := row.Scan(&value)
	if err != nil {
		pgdb.log.Error("Error in get counter value", zap.Error(err))
		return 0, false
	}
	return value, true
}

func (pgdb *PostgresDB) GetAllGaugeValues() map[string]float64 {
	gaugeMetrics := make(map[string]float64)
	var nameValue string
	var value float64
	queryName := `SELECT ID,VALUE FROM metrics WHERE VALUE IS NOT NULL;`
	row, err := pgdb.db.QueryContext(context.Background(),
		queryName)
	if err != nil {
		pgdb.log.Error("Error in get name and gauage value", zap.Error(err))
		return nil
	}
	defer row.Close()
	for row.Next() {
		err := row.Scan(&nameValue, &value)
		if err != nil {
			pgdb.log.Error("Error in get gauage value", zap.Error(err))
			return nil
		}
		gaugeMetrics[nameValue] = value
	}
	return gaugeMetrics
}

func (pgdb *PostgresDB) GetAllCounterValues() map[string]int64 {
	counterMetrics := make(map[string]int64)
	var nameValue string
	var value int64
	queryName := `SELECT ID,DELTA FROM metrics WHERE DELTA IS NOT NULL;`
	row, err := pgdb.db.QueryContext(context.Background(),
		queryName)
	if err != nil {
		pgdb.log.Error("Error in get name and counter value", zap.Error(err))
		return nil
	}
	defer row.Close()
	for row.Next() {
		err := row.Scan(&nameValue, &value)
		if err != nil {
			pgdb.log.Error("Error in get counter value", zap.Error(err))
			return nil
		}
		counterMetrics[nameValue] = value
	}
	return counterMetrics
}
