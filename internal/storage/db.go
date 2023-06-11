package storage

import (
	"context"
	"database/sql"

	"github.com/Zagir2000/alert/internal/models"
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
	tx, err := pgdb.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pgdb *PostgresDB) AddGaugeValue(ctx context.Context, name string, value float64) error {
	tx, err := pgdb.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO metrics (ID,MTYPE,VALUE) VALUES ($1, $2, $3);`, name, "gauge", value)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pgdb *PostgresDB) AddCounterValue(ctx context.Context, name string, value int64) error {
	tx, err := pgdb.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO metrics (ID,MTYPE,DELTA) VALUES ($1, $2, $3);`, name, "counter", value)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pgdb *PostgresDB) GetGauge(ctx context.Context, name string) (float64, bool) {
	var value float64
	row := pgdb.db.QueryRowContext(ctx,
		"SELECT metrics.VALUE  FROM metrics WHERE metrics.ID=$1", name)

	err := row.Scan(&value)
	if err != nil {
		pgdb.log.Error("Error in get gauage value", zap.Error(err))
		return 0, false
	}
	return value, true
}

func (pgdb *PostgresDB) GetCounter(ctx context.Context, name string) (int64, bool) {
	var value int64
	row := pgdb.db.QueryRowContext(ctx,
		"SELECT metrics.DELTA  FROM metrics WHERE metrics.ID=$1", name)
	err := row.Scan(&value)
	if err != nil {
		pgdb.log.Error("Error in get counter value", zap.Error(err))
		return 0, false
	}
	return value, true
}

func (pgdb *PostgresDB) GetAllGaugeValues(ctx context.Context) map[string]float64 {
	gaugeMetrics := make(map[string]float64)
	var nameValue string
	var value float64
	queryName := `SELECT ID,VALUE FROM metrics WHERE VALUE IS NOT NULL;`
	row, err := pgdb.db.QueryContext(ctx,
		queryName)
	lasterr := row.Err()
	if lasterr != nil {
		pgdb.log.Error("Last error in get name and gauge value", zap.Error(err))
	}
	if err != nil {
		pgdb.log.Error("Error in get name and gauge value", zap.Error(err))
		return nil
	}
	defer row.Close()
	for row.Next() {
		err := row.Scan(&nameValue, &value)
		if err != nil {
			pgdb.log.Error("Error in get gauge value", zap.Error(err))
			return nil
		}
		gaugeMetrics[nameValue] = value
	}
	return gaugeMetrics
}

func (pgdb *PostgresDB) GetAllCounterValues(ctx context.Context) map[string]int64 {
	counterMetrics := make(map[string]int64)
	var nameValue string
	var value int64
	queryName := `SELECT ID,DELTA FROM metrics WHERE DELTA IS NOT NULL;`
	row, err := pgdb.db.QueryContext(ctx,
		queryName)
	lasterr := row.Err()
	if lasterr != nil {
		pgdb.log.Error("Last error in get name and counter value", zap.Error(err))
	}
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

func (pgdb *PostgresDB) AddAllValue(ctx context.Context, metrics []models.Metrics) error {
	tx, err := pgdb.db.Begin()
	if err != nil {
		return err
	}
	for _, v := range metrics {
		// все изменения записываются в транзакцию
		if v.MType == "gauge" {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO metrics (ID,MTYPE,VALUE) VALUES ($1, $2, $3);`, v.ID, "gauge", v.Value)
			if err != nil {
				// если ошибка, то откатываем изменения
				tx.Rollback()
				return err
			}
		} else {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO metrics (ID,MTYPE,VALUE) VALUES ($1, $2, $3);`, v.ID, "counter", v.Delta)
			if err != nil {
				// если ошибка, то откатываем изменения
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}
