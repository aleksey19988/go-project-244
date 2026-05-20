package storage

import (
	"encoding/json"
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
	if s.RawData == nil {
		err := s.SetRawData()
		if err != nil {
			return nil, err
		}
	}
	result := make([]map[string]any, 0)
	for _, d := range s.RawData {
		m := make(map[string]any)
		err := json.Unmarshal(d, &m)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}
func (s *Storage) SetRawData() error {
	for _, path := range s.GetPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s.RawData = append(s.RawData, data)
	}

	return nil
}
