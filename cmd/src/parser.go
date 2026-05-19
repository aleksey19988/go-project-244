package parser

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
)

func Parse(s *Storage, format string) (string, error) {
	for _, path := range s.GetPaths() {
		err := validate(path)
		if err != nil {
			return "", err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		s.RawData = append(s.RawData, data)
	}

	mapsWithData, err := s.CreateMapsFromData()
	if err != nil {
		return "", err
	}

	fields, err := diff(mapsWithData)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(fields)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func validate(path string) error {
	if path == "" {
		return errors.New("Path to file is empty")
	}

	ext1 := filepath.Ext(path)
	if ext1 != ".json" {
		return errors.New("Path to file must have .json extension")
	}

	return nil
}

func diff(maps []map[string]any) ([]Field, error) {
	var fields []Field

	keys := make([]string, 0)
	for _, m := range maps {
		for k, _ := range m {
			if !slices.Contains(keys, k) {
				keys = append(keys, k)
			}
		}
	}

	slices.Sort(keys)

	for _, k := range keys {
		isExistsInFirst, isExistsInSecond := false, false
		if _, ok := maps[0][k]; ok {
			isExistsInFirst = true
		}

		if _, ok := maps[1][k]; ok {
			isExistsInSecond = true
		}

		if isExistsInFirst && isExistsInSecond {
			if maps[0][k] == maps[1][k] {
				fields = append(fields, Field{
					Name:     k,
					OldValue: maps[0][k],
					NewValue: maps[1][k],
				})
			} else {
				fields = append(fields, Field{
					Name:     k,
					OldValue: maps[0][k],
					NewValue: nil,
				})
				fields = append(fields, Field{
					Name:     k,
					OldValue: nil,
					NewValue: maps[1][k],
				})
			}
		} else if isExistsInFirst {
			fields = append(fields, Field{
				Name:     k,
				OldValue: maps[0][k],
				NewValue: nil,
			})
		} else if isExistsInSecond {
			fields = append(fields, Field{
				Name:     k,
				OldValue: nil,
				NewValue: maps[1][k],
			})
		}
	}

	return fields, nil
}
