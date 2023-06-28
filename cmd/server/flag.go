package main

import (
	"flag"
	"os"
	"strconv"
)

type FlagVar struct {
	runAddr         string
	logLevel        string
	storeIntervall  int
	fileStoragePath string
	restore         bool
	databaseDsn     string
	migrationsDir   string
	secretKey       string
}

func NewFlagVarStruct() *FlagVar {
	return &FlagVar{}
}

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func (f *FlagVar) parseFlags() error {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию

	flag.StringVar(&f.runAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&f.logLevel, "l", "info", "log level")
	flag.IntVar(&f.storeIntervall, "i", 300, "time interval according to which the current server servers are kept on disk")
	flag.StringVar(&f.fileStoragePath, "f", "/tmp/metrics-db.json", "full name of the file where the current valuee are saved")
	flag.BoolVar(&f.restore, "r", true, "value specifying whether or not to load previously saved values from the specified file at server startup")
	flag.StringVar(&f.databaseDsn, "d", "", "connection to databse")
	flag.StringVar(&f.migrationsDir, "m", "postgresdb/migrations", "migrations to db")
	flag.StringVar(&f.secretKey, "k", "", "key string for signature hash")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.runAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		f.logLevel = envLogLevel
	}
	if envStoreIntervallInt := os.Getenv("STORE_INTERVAL"); envStoreIntervallInt != "" {
		storeIntervallInt, err := strconv.Atoi(envStoreIntervallInt)
		if err != nil {
			return err
		}
		f.storeIntervall = storeIntervallInt
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		f.fileStoragePath = envFileStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		restoreBool, err := strconv.ParseBool(envRestore)
		if err != nil {
			return err
		}
		f.restore = restoreBool
	}
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		f.databaseDsn = envDatabaseDsn
	}
	if envKey, ok := os.LookupEnv("KEY"); ok {
		f.secretKey = envKey
	}
	return nil
}
