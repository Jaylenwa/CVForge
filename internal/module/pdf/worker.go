package pdf

import (
	"context"
	"fmt"
	"time"

	"openresume/internal/infra/cache"
	"openresume/internal/pkg/logger"
	"openresume/internal/pkg/storage"

	"go.uber.org/zap"
)

type Worker struct {
	svc       *Service
	repo      *ExportRepo
	uploader  storage.Uploader
	concur    int
	interval  time.Duration
	timeoutMs int
}

func NewWorker(sys *Service, uploader storage.Uploader) *Worker {
	return &Worker{
		svc:       sys,
		repo:      NewExportRepo(sys),
		uploader:  uploader,
		concur:    4,
		interval:  time.Second,
		timeoutMs: 90000,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	sem := make(chan struct{}, w.concur)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		id, err := w.repo.BLPop(ctx, time.Second*5)
		if err != nil {
			continue
		}
		sem <- struct{}{}
		go func(jobID string) {
			defer func() { <-sem }()
			_ = w.repo.SetProcessing(ctx, jobID)
			job, err := w.repo.GetJob(ctx, jobID)
			if err != nil {
				logger.WithCtx(nil).Error("get job error", zap.Error(err))
				_ = w.repo.SetFailed(ctx, jobID, err.Error())
				return
			}
			buf, code, err := w.svc.GeneratePDFWithToken(job.ResumeID, job.Token)
			if err != nil || code != 200 {
				if err == nil {
					err = fmt.Errorf("pdf code %d", code)
				}
				logger.WithCtx(nil).Error("generate pdf error", zap.Error(err))
				_ = w.repo.SetFailed(ctx, jobID, err.Error())
				return
			}
			name := fmt.Sprintf("resume-%s-%d.pdf", job.ResumeID, time.Now().UnixMilli())
			url, err := w.uploader.Upload(ctx, name, buf)
			if err != nil {
				logger.WithCtx(nil).Error("upload pdf error", zap.Error(err))
				_ = w.repo.SetFailed(ctx, jobID, err.Error())
				return
			}
			_ = cache.RDB.Set(ctx, jobKey(jobID)+":filename", name, time.Hour*24).Err()
			_ = w.repo.SetDone(ctx, jobID, url)
		}(id)
	}
}

func initWorkerDeps(sys *Service) (storage.Uploader, error) {
	cfg := sys.sysConfig.GetStorageSettings()
	return storage.NewFromSettings(cfg)
}

func StartWorker() error {
	svc := NewService()
	uploader, err := initWorkerDeps(svc)
	if err != nil {
		logger.WithCtx(nil).Error("worker deps error", zap.Error(err))
		return err
	}
	w := NewWorker(svc, uploader)
	go func() {
		logger.WithCtx(nil).Info("pdf worker started in-process")
		if err := w.Start(context.Background()); err != nil {
			logger.WithCtx(nil).Error("worker error", zap.Error(err))
		}
	}()
	return nil
}
