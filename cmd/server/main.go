package main

import (
	"log"
	"net/http"

	"github.com/Zagir2000/alert/internal/handlers"
	"github.com/Zagir2000/alert/internal/logger"
	"github.com/Zagir2000/alert/internal/storage"
)

func main() {
	flagStruct := NewFlagVarStruct()
	flagStruct.parseFlags()
	if err := run(flagStruct); err != nil {
		log.Fatalln(err)
	}
}

func run(flagStruct *FlagVar) error {

	log, err := logger.InitializeLogger(flagStruct.logLevel)
	if err != nil {
		return err
	}
	memStorageInterface, posgresDB := storage.NewStorage(log, flagStruct.fileStoragePath, flagStruct.restore, flagStruct.storeIntervall, flagStruct.databaseDsn)
	posgresDB.GetAllGaugeValues()
	defer posgresDB.Close()

	newHandStruct := handlers.MetricHandlerNew(memStorageInterface, log, posgresDB)
	router := handlers.Router(newHandStruct)
	// logger.Log.Info("Running server on", zap.String(flagStruct.runAddr))
	return http.ListenAndServe(flagStruct.runAddr, router)
}
