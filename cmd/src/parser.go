package parser

import (
	"errors"
	"fmt"
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

	formatted := formatOutput(fields, format)

	return formatted, nil
}

func formatOutput(fields []Field, format string) string {
	res := "{\n"

	for _, field := range fields {
		if field.OldValue == field.NewValue {
			res += fmt.Sprintf("    %s: %v\n", field.Name, field.NewValue)
		} else {
			if field.OldValue == nil {
				res += getOutputString(field.TypeOfChange, field.Name, field.NewValue)
			} else {
				res += getOutputString(field.TypeOfChange, field.Name, field.OldValue)
			}
		}
	}

	res += "}"

	return res
}

func getOutputString(typeOfChange, name string, value any) string {
	return fmt.Sprintf("  %s %s: %v\n", typeOfChange, name, value)
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
					Name:         k,
					TypeOfChange: "",
					OldValue:     maps[0][k],
					NewValue:     maps[1][k],
				})
			} else {
				fields = append(fields, Field{
					Name:         k,
					TypeOfChange: Removed,
					OldValue:     maps[0][k],
					NewValue:     nil,
				})
				fields = append(fields, Field{
					Name:         k,
					TypeOfChange: Added,
					OldValue:     nil,
					NewValue:     maps[1][k],
				})
			}
		} else if isExistsInFirst {
			fields = append(fields, Field{
				Name:         k,
				TypeOfChange: Removed,
				OldValue:     maps[0][k],
				NewValue:     nil,
			})
		} else if isExistsInSecond {
			fields = append(fields, Field{
				Name:         k,
				TypeOfChange: Added,
				OldValue:     nil,
				NewValue:     maps[1][k],
			})
		}
	}

	return fields, nil
}
