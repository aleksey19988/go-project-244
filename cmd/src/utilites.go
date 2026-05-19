package parser

import (
	"encoding/json"
)

const (
	Added   = "+"
	Removed = "-"
)

type Storage struct {
	Path1      string
	Path2      string
	RawData    [][]byte
	ParsedData map[string]any
}

type Field struct {
	Name         string
	TypeOfChange string
	OldValue     any
	NewValue     any
}

func (p *Storage) GetPaths() []string {
	return []string{p.Path1, p.Path2}
}
func (p *Storage) CreateMapsFromData() ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	for _, d := range p.RawData {
		m := make(map[string]any)
		err := json.Unmarshal(d, &m)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}
