package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/Zagir2000/alert/internal/logger"
	"go.uber.org/zap"
)

type FlagVar struct {
	runAddr         string
	logLevel        string
	storeIntervall  int
	fileStoragePath string
	restore         bool
}

func NewFlagVarStruct() *FlagVar {
	return &FlagVar{}
}

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func (f *FlagVar) parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&f.runAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&f.logLevel, "l", "info", "log level")
	flag.IntVar(&f.storeIntervall, "i", 300, "time interval according to which the current server servers are kept on disk")
	flag.StringVar(&f.fileStoragePath, "f", "/tmp/metrics-db.json", "full name of the file where the current valuee are saved")
	flag.BoolVar(&f.restore, "r", true, "value specifying whether or not to load previously saved values from the specified file at server startup")

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
			logger.Log.Warn("wrong STORE_INTERVAL format: is not a integer", zap.Error(err))
		}
		f.storeIntervall = storeIntervallInt
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		f.fileStoragePath = envFileStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		restoreBool, err := strconv.ParseBool(envRestore)
		if err != nil {
			logger.Log.Warn("wrong STORE_INTERVAL format: is not a integer", zap.Error(err))
		}
		f.restore = restoreBool
	}
}
