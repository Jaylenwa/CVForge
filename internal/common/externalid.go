package common

import (
	"strings"

	"github.com/google/uuid"
)

func NewExternalID(kind string) string {
	kind = strings.TrimSpace(kind)
	if kind == "" {
		return uuid.NewString()
	}
	return kind + "_" + uuid.NewString()
}
