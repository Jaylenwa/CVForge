package main

import (
	"log"
	"net/http"
	"os"

	"openresume/internal/infra/config"
	"openresume/internal/infra/storage"
	"openresume/internal/module/pdf"
	"openresume/internal/router"
)

func main() {
	config.Load()
	err := storage.Init()
	if err != nil {
		log.Fatalf("storage init error: %v", err)
	}

	if err := pdf.StartInProcessWorker(); err != nil {
		log.Fatalf("pdf worker error: %v", err)
	}

	r := router.Init()
	addr := ":" + config.CF.Port
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
