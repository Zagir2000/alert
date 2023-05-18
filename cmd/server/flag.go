package main

import (
	"flag"
	"os"
)

type FlagVar struct {
	runAddr  string
	logLevel string
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
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.runAddr = envRunAddr
	}
	if envRunAddr := os.Getenv("LOG_LEVEL"); envRunAddr != "" {
		f.logLevel = envRunAddr
	}
}
