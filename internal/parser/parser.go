package parser

import (
	"code/internal/formatters"
	"code/internal/storage"
	"slices"
	"sort"
)

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

	fields, err := Diff(mapsWithData, 1)
	if err != nil {
		return "", err
	}

	if len(format) == 0 {
		format = formatters.DefaultOutputFormat
	}

	f, err := formatters.NewFormatter(format)
	if err != nil {
		return "", err
	}

	formatted, err := f.GetOutput(fields)
	if err != nil {
		return "", err
	}

	return formatted, nil
}

func Diff(pair []map[string]any, deep int) ([]formatters.Field, error) {
	var fields []formatters.Field

	keys := getAllKeys(pair)
	slices.Sort(keys)

	for _, k := range keys {
		firstFieldValue, isKeyExistsInFirst := pair[0][k]
		secondFieldValue, isKeyExistsInSecond := pair[1][k]

		firstMap, firstIsMap := firstFieldValue.(map[string]any)
		secondMap, secondIsMap := secondFieldValue.(map[string]any)

		if isKeyExistsInFirst && isKeyExistsInSecond {
			// ключ есть в обоих файлах
			var children []formatters.Field
			var err error

			if firstIsMap && secondIsMap {
				children, err = Diff([]map[string]any{
					firstMap,
					secondMap,
				}, deep+1)
				if err != nil {
					return nil, err
				}
				fields = append(fields, formatters.Field{
					Name:     k,
					Deep:     deep,
					Children: children,
				})
				continue
			} else if firstIsMap {
				fields = append(fields, formatters.Field{
					Name:         k,
					Deep:         deep,
					TypeOfChange: Removed,
					Children:     getNestedFields(firstMap, deep+1),
				})
				fields = append(fields, formatters.Field{
					Name:         k,
					Deep:         deep,
					TypeOfChange: Added,
					Value:        secondFieldValue,
				})
			} else if secondIsMap {
				fields = append(fields, formatters.Field{
					Name:         k,
					Deep:         deep,
					TypeOfChange: Removed,
					Value:        firstFieldValue,
				})
				fields = append(fields, formatters.Field{
					Name:         k,
					Deep:         deep,
					TypeOfChange: Added,
					Children:     getNestedFields(secondMap, deep+1),
				})
			} else if firstFieldValue != secondFieldValue {
				fields = append(fields, formatters.Field{
					Name:         k,
					TypeOfChange: Removed,
					Value:        firstFieldValue,
					Deep:         deep,
				})
				fields = append(fields, formatters.Field{
					Name:         k,
					TypeOfChange: Added,
					Value:        secondFieldValue,
					Deep:         deep,
				})
			} else {
				fields = append(fields, formatters.Field{
					Name:  k,
					Value: firstFieldValue,
					Deep:  deep,
				})
			}
		} else if isKeyExistsInFirst {
			// ключ есть только в первом, значит был удалён
			fields = append(fields, getRemovedOrAddedField(k, firstFieldValue, Removed, deep))
		} else if isKeyExistsInSecond {
			// ключ есть только во втором, значит был добавлен
			fields = append(fields, getRemovedOrAddedField(k, secondFieldValue, Added, deep))
		}
	}

	return fields, nil
}
func getRemovedOrAddedField(
	fieldName string,
	value any,
	typeOfChange string,
	deep int,
) formatters.Field {
	if mapValue, isMap := value.(map[string]any); isMap {
		return formatters.Field{
			Name:         fieldName,
			TypeOfChange: typeOfChange,
			Deep:         deep,
			Children:     getNestedFields(mapValue, deep+1),
		}
	} else {
		return formatters.Field{
			Name:         fieldName,
			TypeOfChange: typeOfChange,
			Value:        value,
			Deep:         deep,
		}
	}
}
func getNestedFields(m map[string]any, deep int) []formatters.Field {
	var fields []formatters.Field

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		mapValue, ok := m[key].(map[string]any)
		if ok {
			fields = append(fields, formatters.Field{
				Name:     key,
				Children: getNestedFields(mapValue, deep+1),
				Deep:     deep,
			})
		} else {
			fields = append(fields, formatters.Field{
				Name:  key,
				Value: m[key],
				Deep:  deep,
			})
		}
	}

	return fields
}
func getAllKeys(maps []map[string]any) []string {
	res := make([]string, 0)
	keys := map[string]struct{}{}

	for _, m := range maps {
		for k := range m {
			if _, ok := keys[k]; !ok {
				res = append(res, k)
				keys[k] = struct{}{}
			}
		}
	}

	return res
}
