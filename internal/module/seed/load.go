package seed

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

func LoadFromBytes(b []byte) (SeedData, error) {
	var s SeedData
	if err := json.Unmarshal(b, &s); err != nil {
		return SeedData{}, fmt.Errorf("parse seed json: %w", err)
	}
	if err := s.Validate(); err != nil {
		return SeedData{}, err
	}
	return s, nil
}

func LoadFromFS(fsys fs.FS, name string) (SeedData, error) {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		return SeedData{}, fmt.Errorf("read seed file: %w", err)
	}
	return LoadFromBytes(b)
}

func LoadFromFilePath(path string) (SeedData, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return SeedData{}, fmt.Errorf("read seed file: %w", err)
	}
	return LoadFromBytes(b)
}

