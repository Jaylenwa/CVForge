package pdf

import (
	"context"
	"log"

	conf "openresume/internal/module/config"
)

func StartInProcessWorker() error {
	_ = conf.NewService().EnsureDefaults()
	svc := NewService()
	uploader, err := InitWorkerDeps(svc)
	if err != nil {
		log.Printf("worker deps error: %v", err)
		return err
	}
	w := NewWorker(svc, uploader)
	go func() {
		log.Printf("pdf worker started in-process")
		if err := w.Start(context.Background()); err != nil {
			log.Printf("worker error: %v", err)
		}
	}()
	return nil
}
