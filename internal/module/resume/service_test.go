package resume

import (
	"errors"
	"testing"
)

func TestServiceToModelSectionTypeNormalize(t *testing.T) {
	svc := &Service{}
	req := ResumeReq{
		Title:      "t",
		TemplateID: "tpl",
		Language:   "zh",
		Sections: []SectionReq{
			{Type: "Experience", Title: "工作经历", IsVisible: true, Items: []ItemReq{}},
		},
	}
	m, err := svc.toModel(1, req)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(m.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(m.Sections))
	}
	if m.Sections[0].Type != "experience" {
		t.Fatalf("expected normalized type, got %q", m.Sections[0].Type)
	}
}

func TestServiceToModelInvalidSectionType(t *testing.T) {
	svc := &Service{}
	req := ResumeReq{
		Title:      "t",
		TemplateID: "tpl",
		Language:   "zh",
		Sections: []SectionReq{
			{Type: "not_a_real_type", Title: "x", IsVisible: true, Items: []ItemReq{}},
		},
	}
	_, err := svc.toModel(1, req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var inv *InvalidSectionTypeError
	if !errors.As(err, &inv) {
		t.Fatalf("expected InvalidSectionTypeError, got %T", err)
	}
	if inv.Value != "not_a_real_type" {
		t.Fatalf("expected Value preserved, got %q", inv.Value)
	}
}
