package pdf

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"cvforge/internal/infra/cache"
	"cvforge/internal/pkg/logger"
	"cvforge/internal/pkg/storage"

	"go.uber.org/zap"
)

type Worker struct {
	svc       *Service
	repo      *ExportRepo
	uploader  storage.Uploader
	concur    int
	interval  time.Duration
	timeoutMs int
	wg        sync.WaitGroup
}

func NewWorker(sys *Service, uploader storage.Uploader) *Worker {
	wc := 2
	if v := sys.sysConfig.Get("pdf_worker_concurrency"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			wc = n
		}
	}
	return &Worker{
		svc:       sys,
		repo:      NewExportRepo(sys),
		uploader:  uploader,
		concur:    wc,
		interval:  time.Second,
		timeoutMs: 90000,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	sem := make(chan struct{}, w.concur)
	go func() {
		tk := time.NewTicker(30 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				_, _ = w.repo.RequeueStale(ctx, time.Duration(w.timeoutMs)*time.Millisecond)
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			w.wg.Wait()
			return nil
		default:
		}
		id, err := w.repo.BRPopLPush(ctx, time.Second*5)
		if err != nil {
			continue
		}
		sem <- struct{}{}
		go func(jobID string) {
			defer func() {
				<-sem
				w.wg.Done()
			}()
			w.wg.Add(1)
			_ = w.repo.SetProcessing(ctx, jobID)
			job, err := w.repo.GetJob(ctx, jobID)
			if err != nil {
				logger.WithCtx(nil).Error("get job error", zap.Error(err))
				_ = w.repo.SetFailed(ctx, jobID, err.Error())
				_ = w.repo.Confirm(ctx, jobID)
				return
			}
			buf, code, err := w.svc.GeneratePDFWithToken(job.ResumeID, job.Token)
			if err != nil || code != 200 {
				if err == nil {
					err = fmt.Errorf("pdf code %d", code)
				}
				logger.WithCtx(nil).Error("generate pdf error", zap.Error(err))
				_ = w.repo.SetFailed(ctx, jobID, err.Error())
				_ = w.repo.Confirm(ctx, jobID)
				return
			}
			name := fmt.Sprintf("resume-%d-%d.pdf", job.ResumeID, time.Now().UnixMilli())
			url, err := w.uploader.Upload(ctx, name, buf)
			if err != nil {
				logger.WithCtx(nil).Error("upload pdf error", zap.Error(err))
				_ = w.repo.SetFailed(ctx, jobID, err.Error())
				_ = w.repo.Confirm(ctx, jobID)
				return
			}
			_ = cache.RDB.Set(ctx, jobKey(jobID)+":filename", name, time.Hour*24).Err()
			_ = w.repo.SetDone(ctx, jobID, url)
			_ = w.repo.Confirm(ctx, jobID)
		}(id)
	}
}

func initWorkerDeps(sys *Service) (storage.Uploader, error) {
	cfg := sys.sysConfig.GetStorageSettings()
	return storage.NewFromSettings(cfg)
}

func StartWorker(ctx context.Context) error {
	svc := NewService()
	uploader, err := initWorkerDeps(svc)
	if err != nil {
		logger.WithCtx(nil).Error("worker deps error", zap.Error(err))
		return err
	}
	w := NewWorker(svc, uploader)
	go func() {
		logger.WithCtx(nil).Info("pdf worker started in-process")
		if err := w.Start(ctx); err != nil {
			logger.WithCtx(nil).Error("worker error", zap.Error(err))
		}
	}()
	return nil
}
