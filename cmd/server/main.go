package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Zagir2000/alert/internal/server/handlers"
	"github.com/Zagir2000/alert/internal/server/logger"
	"github.com/Zagir2000/alert/internal/server/storage"
	"go.uber.org/zap"
)

func main() {
	flagStruct := NewFlagVarStruct()
	err := flagStruct.parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	if err := run(flagStruct); err != nil {
		log.Fatalln(err)
	}
}

func run(flagStruct *FlagVar) error {

	log, err := logger.InitializeLogger(flagStruct.logLevel)
	if err != nil {
		return err
	}
	ctx := context.Background()
	memStorageInterface, postgresDB, err := storage.NewStorage(ctx, flagStruct.migrationsDir, log, flagStruct.fileStoragePath, flagStruct.restore, flagStruct.storeIntervall, flagStruct.databaseDsn)
	if err != nil {
		log.Fatal("Error in create storage", zap.Error(err))
	}
	if postgresDB != nil {
		defer postgresDB.Close()
	}

	newHandStruct := handlers.MetricHandlerNew(memStorageInterface, log, postgresDB)
	router := handlers.Router(ctx, newHandStruct)
	log.Info("Running server on", zap.String("", flagStruct.runAddr))
	return http.ListenAndServe(flagStruct.runAddr, router)
}
