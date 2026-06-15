package parser

import (
	"code/internal/diff"
	"code/internal/formatters"
	"slices"
	"sort"
)

func Diff(pair []map[string]any, deep int) ([]diff.Field, error) {
	var fields []diff.Field

	keys := getAllKeys(pair)
	slices.Sort(keys)

	for _, k := range keys {
		firstFieldValue, isKeyExistsInFirst := pair[0][k]
		secondFieldValue, isKeyExistsInSecond := pair[1][k]

		firstMap, firstIsMap := firstFieldValue.(map[string]any)
		secondMap, secondIsMap := secondFieldValue.(map[string]any)

		if isKeyExistsInFirst && isKeyExistsInSecond {
			// ключ есть в обоих файлах
			var children []diff.Field
			var err error

			if firstIsMap && secondIsMap {
				children, err = Diff([]map[string]any{
					firstMap,
					secondMap,
				}, deep+1)
				if err != nil {
					return nil, err
				}
				fields = append(fields, diff.Field{
					Name:     k,
					Depth:    deep,
					Children: children,
				})
				continue
			} else if firstIsMap {
				fields = append(fields, diff.Field{
					Name:     k,
					Depth:    deep,
					Status:   formatters.Updated,
					OldValue: getNestedFields(firstMap, deep+1),
					NewValue: secondFieldValue,
					Children: getNestedFields(firstMap, deep+1),
				})
			} else if secondIsMap {
				fields = append(fields, diff.Field{
					Name:     k,
					Depth:    deep,
					Status:   formatters.Updated,
					OldValue: firstFieldValue,
					NewValue: getNestedFields(secondMap, deep+1),
					Children: getNestedFields(secondMap, deep+1),
				})
			} else if firstFieldValue != secondFieldValue {
				fields = append(fields, diff.Field{
					Name:     k,
					Depth:    deep,
					Status:   formatters.Updated,
					OldValue: firstFieldValue,
					NewValue: secondFieldValue,
				})
			} else {
				fields = append(fields, diff.Field{
					Name:     k,
					OldValue: firstFieldValue,
					NewValue: firstFieldValue,
					Status:   formatters.Unchanged,
					Depth:    deep,
				})
			}
		} else if isKeyExistsInFirst {
			// ключ есть только в первом, значит был удалён
			fields = append(fields, getRemovedOrAddedField(k, firstFieldValue, formatters.Removed, deep))
		} else if isKeyExistsInSecond {
			// ключ есть только во втором, значит был добавлен
			fields = append(fields, getRemovedOrAddedField(k, secondFieldValue, formatters.Added, deep))
		}
	}

	return fields, nil
}
func getRemovedOrAddedField(
	fieldName string,
	value any,
	typeOfChange string,
	deep int,
) diff.Field {
	result := diff.Field{
		Name:   fieldName,
		Status: typeOfChange,
		Depth:  deep,
	}
	if mapValue, isMap := value.(map[string]any); isMap {
		result.Children = getNestedFields(mapValue, deep+1)
	} else {
		result.OldValue = value
		result.NewValue = value
	}

	return result
}
func getNestedFields(m map[string]any, deep int) []diff.Field {
	var fields []diff.Field

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		mapValue, ok := m[key].(map[string]any)
		if ok {
			fields = append(fields, diff.Field{
				Name:     key,
				Children: getNestedFields(mapValue, deep+1),
				Depth:    deep,
			})
		} else {
			fields = append(fields, diff.Field{
				Name:     key,
				OldValue: m[key],
				NewValue: m[key],
				Depth:    deep,
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
