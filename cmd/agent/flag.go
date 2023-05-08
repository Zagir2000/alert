package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func parseFlags() (string, int, int) {
	var runAddr string
	var reportInterval int
	var pollInterval int
	// как аргумент -a со значением :8080 по умолчанию
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")

	// частота отправки метрик на сервер
	flag.IntVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")

	//частота опроса метрик из пакета
	flag.IntVar(&pollInterval, "p", 2, "frequency of polling metrics from the package")
	flag.Parse()
	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		runAddr = envRunAddr
	}
	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		envReportIntervalInt, err := strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatalln("wrong REPORT_INTERVAL format: is not a integer", err)
		}
		reportInterval = envReportIntervalInt
	}
	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		envPollIntervalInt, err := strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatalln("wrong POLL_INTERVAL format: is not a integer", err)
		}
		pollInterval = envPollIntervalInt
	}
	return runAddr, pollInterval, reportInterval
}
