package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	YamlExt = ".yaml"
	YmlExt  = ".yml"
	JsonExt = ".json"
)

type Storage struct {
	Path1     string
	Path2     string
	RawData   [][]byte
	Extension string
}

func NewStorage(path1, path2 string) (*Storage, error) {
	extensions := map[string]int{}

	for key, path := range []string{path1, path2} {
		if path == "" {
			return nil, fmt.Errorf("path to file %d is empty", key+1)
		}

		ext := filepath.Ext(path)
		if ext != JsonExt && ext != YamlExt && ext != YmlExt {
			return nil, fmt.Errorf("path to file %d must have .json or .yaml extension", key+1)
		}

		extensions[ext]++
	}

	if len(extensions) > 1 {
		return nil, errors.New("files must have one extension")
	}

	return &Storage{
		Path1:     path1,
		Path2:     path2,
		Extension: filepath.Ext(path1),
	}, nil
}
func (s *Storage) GetPaths() []string {
	return []string{s.Path1, s.Path2}
}
func (s *Storage) CreateMapsFromData() ([]map[string]any, error) {
	if len(s.RawData) == 0 {
		if err := s.SetRawData(); err != nil {
			return nil, err
		}
	}
	result := make([]map[string]any, 0)
	for i, data := range s.RawData {
		m := make(map[string]any)
		switch s.Extension {
		case JsonExt:
			if err := json.Unmarshal(data, &m); err != nil {
				return nil, fmt.Errorf("failed to unmarshal file #%d: %w", i+1, err)
			}
		case YamlExt, YmlExt:
			if err := yaml.Unmarshal(data, &m); err != nil {
				return nil, fmt.Errorf("failed to unmarshal file #%d: %w", i+1, err)
			}
		}
		result = append(result, m)
	}

	return result, nil
}
func (s *Storage) SetRawData() error {
	s.RawData = make([][]byte, 0, 2)

	for _, path := range s.GetPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}
		s.RawData = append(s.RawData, data)
	}
	return nil
}
