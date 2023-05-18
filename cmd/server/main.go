package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/server/handlers"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
func run() error {
	flagStruct := NewFlagVarStruct()
	flagStruct.parseFlags()
	err := logger.Initialize(flagStruct.logLevel)
	if err != nil {
		return err
	}
	router := handlers.Router()
	fmt.Println("Running server on", flagStruct.runAddr)
	return http.ListenAndServe(flagStruct.runAddr, router)
}
