package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zagir2000/alert/internal/server/handlers"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
func run() error {
	router := handlers.Router()
	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, router)
}
