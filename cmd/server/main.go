package main

import (
	"context"
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
	if flagStruct.restore {
		err := memStorage.MetricsLoadJSON(flagStruct.fileStoragePath)
		if err != nil {
			logger.Log.Error("failed to load file", zap.Error(err))
		}
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		go func() {
			select {
			case <-ctx.Done():
				return
			default:
				for {
					err = storage.MetricsSaveJSON(flagStruct.fileStoragePath, memStorage)
					if err != nil {
						logger.Log.Error("failed to save file", zap.Error(err))
						cancel()
					}
					time.Sleep(time.Duration(flagStruct.storeIntervall) * time.Second)
				}
			}
		}()

	}

	newHandStruct := handlers.MetricHandlerNew(memStorage)
	router := handlers.Router(newHandStruct)
	// logger.Log.Info("Running server on", zap.String(flagStruct.runAddr))
	return http.ListenAndServe(flagStruct.runAddr, router)
}
