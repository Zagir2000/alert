package main

import (
	"flag"
	"os"
	"strconv"
)

type FlagVar struct {
	runAddr        string
	reportInterval int
	pollInterval   int
}

func NewFlagVarStruct() *FlagVar {
	return &FlagVar{}
}
func (f *FlagVar) parseFlags() error {

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
			return err
		}
		f.reportInterval = envReportIntervalInt
	}
	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		envPollIntervalInt, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return err
		}
		f.pollInterval = envPollIntervalInt
	}
	return nil
}
