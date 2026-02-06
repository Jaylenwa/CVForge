package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"openresume/internal/infra/config"
	"openresume/internal/infra/storage"
	"openresume/internal/module/pdf"
	"openresume/internal/module/seed"
	"openresume/internal/pkg/logger"
	"openresume/internal/router"

	"go.uber.org/zap"
)

func main() {
	logger.Init()
	config.Load()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := storage.Init(); err != nil {
			logger.L().Fatal("storage init error", zap.Error(err))
		}
		os.Exit(seed.RunCLI(os.Args[2:]))
	}

	err := storage.Init()
	if err != nil {
		logger.L().Fatal("storage init error", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := pdf.StartWorker(ctx); err != nil {
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
