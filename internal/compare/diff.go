package compare

import (
	"slices"
	"sort"
)

const (
	Added     = "added"
	Removed   = "removed"
	Updated   = "updated"
	Unchanged = "unchanged"
)

func GenDiff(first, second map[string]any) []Field {
	return genDiffWithDepth(first, second, 1)
}
func genDiffWithDepth(first, second map[string]any, depth int) []Field {
	var fields []Field

	keys := getAllKeys([]map[string]any{first, second})
	slices.Sort(keys)

	for _, k := range keys {
		firstFieldValue, isKeyExistsInFirst := first[k]
		secondFieldValue, isKeyExistsInSecond := second[k]

		firstMap, firstIsMap := firstFieldValue.(map[string]any)
		secondMap, secondIsMap := secondFieldValue.(map[string]any)

		if isKeyExistsInFirst && isKeyExistsInSecond {
			// ключ есть в обоих файлах
			var children []Field

			if firstIsMap && secondIsMap {
				children = genDiffWithDepth(firstMap, secondMap, depth+1)
				fields = append(fields, Field{
					Name:     k,
					Depth:    depth,
					Children: children,
				})
				continue
			} else if firstIsMap {
				fields = append(fields, Field{
					Name:     k,
					Depth:    depth,
					Status:   Updated,
					OldValue: getNestedFields(firstMap, depth+1),
					NewValue: secondFieldValue,
					Children: getNestedFields(firstMap, depth+1),
				})
			} else if secondIsMap {
				fields = append(fields, Field{
					Name:     k,
					Depth:    depth,
					Status:   Updated,
					OldValue: firstFieldValue,
					NewValue: getNestedFields(secondMap, depth+1),
					Children: getNestedFields(secondMap, depth+1),
				})
			} else if firstFieldValue != secondFieldValue {
				fields = append(fields, Field{
					Name:     k,
					Depth:    depth,
					Status:   Updated,
					OldValue: firstFieldValue,
					NewValue: secondFieldValue,
				})
			} else {
				fields = append(fields, Field{
					Name:     k,
					OldValue: firstFieldValue,
					NewValue: firstFieldValue,
					Status:   Unchanged,
					Depth:    depth,
				})
			}
		} else if isKeyExistsInFirst {
			// ключ есть только в первом, значит был удалён
			fields = append(fields, getRemovedOrAddedField(k, firstFieldValue, Removed, depth))
		} else if isKeyExistsInSecond {
			// ключ есть только во втором, значит был добавлен
			fields = append(fields, getRemovedOrAddedField(k, secondFieldValue, Added, depth))
		}
	}

	return fields
}
func getRemovedOrAddedField(
	fieldName string,
	value any,
	typeOfChange string,
	depth int,
) Field {
	result := Field{
		Name:   fieldName,
		Status: typeOfChange,
		Depth:  depth,
	}
	if mapValue, isMap := value.(map[string]any); isMap {
		result.Children = getNestedFields(mapValue, depth+1)
	} else {
		result.OldValue = value
		result.NewValue = value
	}

	return result
}
func getNestedFields(m map[string]any, depth int) []Field {
	var fields []Field

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		mapValue, ok := m[key].(map[string]any)
		if ok {
			fields = append(fields, Field{
				Name:     key,
				Children: getNestedFields(mapValue, depth+1),
				Depth:    depth,
			})
		} else {
			fields = append(fields, Field{
				Name:     key,
				OldValue: m[key],
				NewValue: m[key],
				Depth:    depth,
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
