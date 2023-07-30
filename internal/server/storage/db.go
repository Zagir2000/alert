package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Zagir2000/alert/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type PostgresDB struct {
	pool *pgxpool.Pool
	log  *zap.Logger
	rw   sync.RWMutex
}

func (pgdb *PostgresDB) PingDB(ctx context.Context) error {
	err := pgdb.pool.Ping(ctx)
	return err
}

func InitDB(configDB string, log *zap.Logger, migratePath string) (*PostgresDB, error) {
	err := runMigrations(configDB, migratePath)
	if err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), configDB)
	if err == nil {
		return &PostgresDB{pool: pool, log: log}, nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		log.Error("Database initialization error", zap.Error(err))
		for _, k := range models.TimeConnect {
			time.Sleep(k)
			pool, err := pgxpool.New(context.Background(), configDB)
			if err == nil {
				log.Info("Successful database connection")
				return &PostgresDB{pool: pool, log: log}, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to create a connection pool: %w", err)
}

func runMigrations(dsn string, migratePath string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migratePath), dsn)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (pgdb *PostgresDB) Close() {
	pgdb.pool.Close()
}

func (pgdb *PostgresDB) AddGaugeValue(ctx context.Context, name string, value float64) error {
	pgdb.rw.Lock()
	defer pgdb.rw.Unlock()
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO metrics (mname,mtype,value) VALUES ($1, $2, $3) ON CONFLICT (mname, mtype) DO UPDATE SET value = $3;`, name, "gauge", value)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (pgdb *PostgresDB) AddCounterValue(ctx context.Context, name string, value int64) error {
	pgdb.rw.Lock()
	defer pgdb.rw.Unlock()
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO metrics (mname,mtype,delta) VALUES ($1, $2, $3)  ON CONFLICT (mname, mtype) DO UPDATE SET delta = metrics.delta+$3;`, name, "counter", value)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (pgdb *PostgresDB) GetGauge(ctx context.Context, name string) (float64, bool) {
	var value float64
	row := pgdb.pool.QueryRow(ctx,
		"SELECT metrics.value  FROM metrics WHERE metrics.mname=$1", name)

	err := row.Scan(&value)
	if err != nil {
		errStr := fmt.Sprintf("Error in get gauge value %s", name)
		pgdb.log.Error(errStr, zap.Error(err))
		return 0, false
	}
	return value, true
}

func (pgdb *PostgresDB) GetCounter(ctx context.Context, name string) (int64, bool) {
	var value int64
	row := pgdb.pool.QueryRow(ctx,
		"SELECT metrics.delta  FROM metrics WHERE metrics.mname=$1", name)
	err := row.Scan(&value)
	if err != nil {
		errStr := fmt.Sprintf("Error in get counter value %s", name)
		pgdb.log.Error(errStr, zap.Error(err))
		return 0, false
	}
	return value, true
}

func (pgdb *PostgresDB) GetAllGaugeValues(ctx context.Context) map[string]float64 {
	gaugeMetrics := make(map[string]float64)
	var nameValue string
	var value float64
	queryName := `SELECT mname,value FROM metrics WHERE value IS NOT NULL;`
	row, err := pgdb.pool.Query(ctx,
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
	queryName := `SELECT manme,delta FROM metrics WHERE delta IS NOT NULL;`
	row, err := pgdb.pool.Query(ctx,
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
	pgdb.rw.Lock()
	defer pgdb.rw.Unlock()
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {
		return err
	}
	for _, v := range metrics {
		// все изменения записываются в транзакцию
		if v.MType == "gauge" {
			_, err = tx.Exec(ctx,
				`INSERT INTO metrics (mname,mtype,value) VALUES ($1, $2, $3) ON CONFLICT (mname, mtype) DO UPDATE SET value = $3;`, v.ID, "gauge", v.Value)
			if err != nil {
				// если ошибка, то откатываем изменения
				tx.Rollback(ctx)
				return err
			}
		} else {
			_, err = tx.Exec(ctx,
				`INSERT INTO metrics (mname,mtype,delta) VALUES ($1, $2, $3)  ON CONFLICT (mname, mtype) DO UPDATE SET delta = metrics.delta+$3;`, v.ID, "counter", v.Delta)
			if err != nil {
				// если ошибка, то откатываем изменения
				tx.Rollback(ctx)
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
