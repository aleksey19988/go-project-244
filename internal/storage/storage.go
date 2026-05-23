package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storage struct {
	Path1   string
	Path2   string
	RawData [][]byte
}

func NewStorage(path1, path2 string) *Storage {
	return &Storage{Path1: path1, Path2: path2}
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
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file #%d: %w", i+1, err)
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
