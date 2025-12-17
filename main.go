package main

import (
	"log"
	"net/http"
	"os"

	"openresume/internal/infra/config"
	"openresume/internal/infra/storage"
	"openresume/internal/router"
)

func main() {
	cfg := config.Load()

	db, rdb, err := storage.Init(cfg)
	if err != nil {
		log.Fatalf("storage init error: %v", err)
	}

	r := router.Init(cfg, db, rdb)

	addr := ":" + cfg.Port
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
