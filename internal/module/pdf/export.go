package pdf

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"openresume/internal/infra/cache"

	"github.com/google/uuid"
)

type ExportStatus string

const (
	ExportStatusPending    ExportStatus = "pending"
	ExportStatusProcessing ExportStatus = "processing"
	ExportStatusDone       ExportStatus = "done"
	ExportStatusFailed     ExportStatus = "failed"
)

type ExportJob struct {
	ID        string            `json:"id"`
	UserID    string            `json:"userId"`
	ResumeID  uint              `json:"resumeId"`
	Token     string            `json:"token"`
	Options   map[string]string `json:"options,omitempty"`
	CreatedAt int64             `json:"createdAt"`
}

func queueKey() string        { return "queue:pdf" }
func jobKey(id string) string { return "job:pdf:" + id }

type ExportRepo struct {
	sysConfig *Service
}

func NewExportRepo(sys *Service) *ExportRepo {
	return &ExportRepo{sysConfig: sys}
}

func (r *ExportRepo) Enqueue(ctx context.Context, job ExportJob) error {
	if cache.RDB == nil {
		return fmt.Errorf("redis not initialized")
	}
	job.CreatedAt = time.Now().Unix()
	b, err := json.Marshal(job)
	if err != nil {
		return err
	}
	maxQLen := int64(1000)
	if v := r.sysConfig.sysConfig.Get("pdf_queue_max"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			maxQLen = n
		}
	}
	if qlen := cache.RDB.LLen(ctx, queueKey()).Val(); qlen >= maxQLen {
		return fmt.Errorf("queue full")
	}
	if err := cache.RDB.Set(ctx, jobKey(job.ID)+":status", string(ExportStatusPending), time.Hour*24).Err(); err != nil {
		return err
	}
	if err := cache.RDB.Set(ctx, jobKey(job.ID)+":data", b, time.Hour*24).Err(); err != nil {
		return err
	}
	return cache.RDB.RPush(ctx, queueKey(), job.ID).Err()
}

func (r *ExportRepo) SetProcessing(ctx context.Context, id string) error {
	return cache.RDB.Set(ctx, jobKey(id)+":status", string(ExportStatusProcessing), time.Hour*24).Err()
}

func (r *ExportRepo) SetDone(ctx context.Context, id string, url string) error {
	if err := cache.RDB.Set(ctx, jobKey(id)+":status", string(ExportStatusDone), time.Hour*24).Err(); err != nil {
		return err
	}
	return cache.RDB.Set(ctx, jobKey(id)+":result", url, time.Hour*24).Err()
}

func (r *ExportRepo) SetFailed(ctx context.Context, id string, msg string) error {
	if err := cache.RDB.Set(ctx, jobKey(id)+":status", string(ExportStatusFailed), time.Hour*24).Err(); err != nil {
		return err
	}
	return cache.RDB.Set(ctx, jobKey(id)+":error", msg, time.Hour*24).Err()
}

func (r *ExportRepo) GetStatus(ctx context.Context, id string) (ExportStatus, string, string, error) {
	st := cache.RDB.Get(ctx, jobKey(id)+":status").Val()
	res := cache.RDB.Get(ctx, jobKey(id)+":result").Val()
	errMsg := cache.RDB.Get(ctx, jobKey(id)+":error").Val()
	if st == "" {
		return "", "", "", fmt.Errorf("not found")
	}
	return ExportStatus(st), res, errMsg, nil
}

func (r *ExportRepo) BLPop(ctx context.Context, timeout time.Duration) (string, error) {
	val, err := cache.RDB.BLPop(ctx, timeout, queueKey()).Result()
	if err != nil {
		return "", err
	}
	if len(val) != 2 {
		return "", fmt.Errorf("invalid blpop result")
	}
	return val[1], nil
}

func (r *ExportRepo) GetJob(ctx context.Context, id string) (ExportJob, error) {
	var job ExportJob
	raw := cache.RDB.Get(ctx, jobKey(id)+":data").Val()
	if raw == "" {
		return job, fmt.Errorf("not found")
	}
	if err := json.Unmarshal([]byte(raw), &job); err != nil {
		return job, err
	}
	return job, nil
}

func (r *ExportRepo) EnsureOTT(ctx context.Context, id string) (string, error) {
	k := jobKey(id) + ":ott"
	v := cache.RDB.Get(ctx, k).Val()
	if v != "" {
		return v, nil
	}
	t := uuid.NewString()
	ok := cache.RDB.SetNX(ctx, k, t, 10*time.Minute).Val()
	if !ok {
		v = cache.RDB.Get(ctx, k).Val()
		if v == "" {
			return "", fmt.Errorf("ott unavailable")
		}
		return v, nil
	}
	return t, nil
}

func (r *ExportRepo) ValidateOTT(ctx context.Context, id string, token string) bool {
	k := jobKey(id) + ":ott"
	v := cache.RDB.Get(ctx, k).Val()
	if v == "" || v != token {
		return false
	}
	_ = cache.RDB.Del(ctx, k).Err()
	return true
}
