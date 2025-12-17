package ai

import "strings"

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Polish(text, tone string) string {
	out := strings.TrimSpace(text)
	if out != "" {
		out = strings.ToUpper(out[:1]) + out[1:]
	}
	return out
}

func (s *Service) Summary(jobTitle, skills string) string {
	return "Experienced " + jobTitle + " with skills in " + skills + ". Delivers impact through collaboration and ownership."
}

