package template

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
)

type TemplateDTO struct {
	ExternalID string
	Name       string
	Names      map[string]string `json:",omitempty"`
	UsageCount int
}

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func pickName(names map[string]string, language string) string {
	if len(names) == 0 {
		return ""
	}
	if v := strings.TrimSpace(names[language]); v != "" {
		return v
	}
	if v := strings.TrimSpace(names["zh"]); v != "" {
		return v
	}
	for _, v := range names {
		if v = strings.TrimSpace(v); v != "" {
			return v
		}
	}
	return ""
}

func (s *Service) ListAllPayload() (string, error) {
	if val, err := cache.RDB.Get(context.Background(), string(common.RedisKeyTemplatesListAll)).Result(); err == nil {
		return val, nil
	}
	list, err := s.repo.ListAll()
	if err != nil {
		return "", err
	}

	ids := make([]uint, 0, len(list))
	for _, it := range list {
		ids = append(ids, it.ID)
	}
	i18nList, err := s.repo.ListI18n(ids)
	if err != nil {
		return "", err
	}
	namesByID := map[uint]map[string]string{}
	for _, it := range i18nList {
		m := namesByID[it.TemplateID]
		if m == nil {
			m = map[string]string{}
			namesByID[it.TemplateID] = m
		}
		m[it.Language] = it.Name
	}
	out := make([]TemplateDTO, 0, len(list))
	for _, it := range list {
		names := namesByID[it.ID]
		name := pickName(names, "zh")
		out = append(out, TemplateDTO{
			ExternalID: it.ExternalID,
			Name:       name,
			Names:      names,
			UsageCount: it.UsageCount,
		})
	}

	payloadBytes, _ := json.Marshal(map[string]any{"items": out})
	payload := string(payloadBytes)
	cache.RDB.Set(context.Background(), string(common.RedisKeyTemplatesListAll), payload, time.Hour)
	return payload, nil
}

func (s *Service) GetByExternal(id string) (TemplateDTO, error) {
	t, err := s.repo.GetByExternal(id)
	if err != nil {
		return TemplateDTO{}, err
	}
	i18nList, err := s.repo.ListI18n([]uint{t.ID})
	if err != nil {
		return TemplateDTO{}, err
	}
	names := map[string]string{}
	for _, it := range i18nList {
		names[it.Language] = it.Name
	}
	name := pickName(names, "zh")
	return TemplateDTO{
		ExternalID: t.ExternalID,
		Name:       name,
		Names:      names,
		UsageCount: t.UsageCount,
	}, nil
}

type SeedTemplateItem struct {
	ExternalID string
	Names      map[string]string
}

func (s *Service) Seed(items []SeedTemplateItem) error {
	for _, it := range items {
		externalID := strings.TrimSpace(it.ExternalID)
		if externalID == "" {
			continue
		}
		t, err := s.repo.GetByExternal(externalID)
		if err != nil {
			created := Template{ExternalID: externalID}
			if err := s.repo.CreateWithNames(&created, it.Names); err != nil {
				return err
			}
			continue
		}
		if err := s.repo.PatchWithNames(t.ID, nil, it.Names); err != nil {
			return err
		}
	}
	if cache.RDB != nil {
		cache.RDB.Del(context.Background(), string(common.RedisKeyTemplatesListAll))
	}
	return nil
}
