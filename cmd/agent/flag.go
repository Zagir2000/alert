package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/Zagir2000/alert/internal/logger"
	"go.uber.org/zap"
)

type FlagVar struct {
	runAddr        string
	reportInterval int
	pollInterval   int
}

func NewFlagVarStruct() *FlagVar {
	return &FlagVar{}
}
func (f *FlagVar) parseFlags() {

	// как аргумент -a со значением :8080 по умолчанию
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.StringVar(&f.runAddr, "a", "localhost:8080", "address and port to run server")

	// частота отправки метрик на сервер
	flag.IntVar(&f.reportInterval, "r", 10, "frequency of sending metrics to the server")

	//частота опроса метрик из пакета
	flag.IntVar(&f.pollInterval, "p", 2, "frequency of polling metrics from the package")
	flag.Parse()
	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		f.runAddr = envRunAddr
	}
	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		envReportIntervalInt, err := strconv.Atoi(envReportInterval)
		if err != nil {
			logger.Log.Warn("wrong REPORT_INTERVAL format: is not a integer", zap.Error(err))
		}
		f.reportInterval = envReportIntervalInt
	}
	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		envPollIntervalInt, err := strconv.Atoi(envPollInterval)
		if err != nil {
			logger.Log.Warn("wrong POLL_INTERVAL format: is not a integer", zap.Error(err))
		}
		f.pollInterval = envPollIntervalInt
	}
}
