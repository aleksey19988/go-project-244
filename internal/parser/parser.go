package parser

import (
	"code/internal/storage"
	"fmt"
	"slices"
)

type Field struct {
	Name         string
	TypeOfChange string
	OldValue     any
	NewValue     any
}

const (
	Added   = "+"
	Removed = "-"
)

func Parse(s *storage.Storage, format string) (string, error) {
	err := s.SetRawData()
	if err != nil {
		return "", err
	}

	mapsWithData, err := s.CreateMapsFromData()
	if err != nil {
		return "", err
	}

	fields, err := Diff(mapsWithData)
	if err != nil {
		return "", err
	}

	formatted := FormatOutput(fields, format)

	return formatted, nil
}

func FormatOutput(fields []Field, format string) string {
	res := "{\n"

	for _, field := range fields {
		if field.OldValue == field.NewValue {
			res += fmt.Sprintf("    %s: %v\n", field.Name, field.NewValue)
		} else {
			if field.OldValue == nil {
				res += GetOutputString(field.TypeOfChange, field.Name, field.NewValue)
			} else {
				res += GetOutputString(field.TypeOfChange, field.Name, field.OldValue)
			}
		}
	}

	res += "}"

	return res
}
func GetOutputString(typeOfChange, name string, value any) string {
	return fmt.Sprintf("  %s %s: %v\n", typeOfChange, name, value)
}
func Diff(maps []map[string]any) ([]Field, error) {
	var fields []Field

	keys := getCommonKeys(maps)
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
func getCommonKeys(maps []map[string]any) []string {
	res := make([]string, 0)

	for _, m := range maps {
		for k := range m {
			if !slices.Contains(res, k) {
				res = append(res, k)
			}
		}
	}

	return res
}
