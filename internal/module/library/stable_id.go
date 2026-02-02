package library

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func stableExternalID(kind string, parts ...string) string {
	h := sha1.New()
	h.Write([]byte(strings.TrimSpace(kind)))
	for _, p := range parts {
		h.Write([]byte{0})
		h.Write([]byte(strings.TrimSpace(p)))
	}
	sum := hex.EncodeToString(h.Sum(nil))
	return kind + "_" + sum[:16]
}
