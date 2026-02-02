package main

import (
	"net/http"
	"os"

	"go.uber.org/zap"
	"openresume/internal/infra/config"
	"openresume/internal/infra/storage"
	"openresume/internal/module/pdf"
	"openresume/internal/pkg/logger"
	"openresume/internal/router"
)

func main() {
	logger.Init()
	config.Load()
	err := storage.Init()
	if err != nil {
		logger.L().Fatal("storage init error", zap.Error(err))
	}

	if err := pdf.StartWorker(); err != nil {
		logger.L().Fatal("pdf worker error", zap.Error(err))
	}

	r := router.Init()
	addr := ":" + config.CF.Port
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}

	logger.L().Info("server listening", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.L().Fatal("server error", zap.Error(err))
	}
}
