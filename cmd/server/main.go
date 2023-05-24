package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zagir2000/alert/internal/handlers"
	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/storage"
	"go.uber.org/zap"
)

func main() {
	flagStruct := NewFlagVarStruct()
	flagStruct.parseFlags()
	if err := run(flagStruct); err != nil {
		log.Fatalln(err)
	}
}
func run(flagStruct *FlagVar) error {

	err := logger.Initialize(flagStruct.logLevel)
	if err != nil {
		return err
	}
	memStorage := storage.NewMemStorage()
	if flagStruct.restore == true {
		err := storage.MetricsLoadJSON(flagStruct.fileStoragePath, memStorage)
		if err != nil {
			logger.Log.Error("failed to load file", zap.Error(err))
		}
		go func() {
			for {
				time.Sleep(time.Duration(flagStruct.storeIntervall) * time.Second)
				err = storage.MetricsSaveJson(flagStruct.fileStoragePath, memStorage)
				if err != nil {
					logger.Log.Error("failed to save file", zap.Error(err))
				}
			}
		}()

	}

	newHandStruct := handlers.MetricHandlerNew(memStorage)
	router := handlers.Router(newHandStruct)
	fmt.Println("Running server on", flagStruct.runAddr)
	return http.ListenAndServe(flagStruct.runAddr, router)
}
