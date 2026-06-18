package formatters

import (
	"code/internal/compare"
	"encoding/json"
)

type JsonFormatter struct{}

func (jf *JsonFormatter) Format(fields []compare.Field) (string, error) {
	res, err := getFieldDataForOutput(fields)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func getFieldDataForOutput(fds []compare.Field) (map[string]any, error) {
	res := make(map[string]any)

	for _, field := range fds {
		fieldData := map[string]any{}

		if field.Status != "" {
			fieldData["status"] = field.Status
		}

		switch field.Status {
		case compare.Added:
			value, err := parseFieldForJson(field, compare.Added)
			if err != nil {
				return nil, err
			}
			fieldData["value"] = value
			res[field.Name] = fieldData
		case compare.Removed, compare.Unchanged:
			value, err := parseFieldForJson(field, compare.Removed)
			if err != nil {
				return nil, err
			}
			fieldData["value"] = value
			res[field.Name] = fieldData
		case compare.Updated:
			oldValue, err := parseFieldForJson(field, compare.Removed)
			if err != nil {
				return nil, err
			}
			fieldData["old_value"] = oldValue

			newValue, err := parseFieldForJson(field, compare.Added)
			if err != nil {
				return nil, err
			}
			fieldData["new_value"] = newValue
			res[field.Name] = fieldData
		default:
			if len(field.Children) > 0 {
				childrenData, err := getFieldDataForOutput(field.Children)
				if err != nil {
					return nil, err
				}
				for k, v := range childrenData {
					fieldData[k] = v
				}
				res[field.Name] = fieldData
			} else {
				res[field.Name] = field.OldValue
			}
		}
	}

	return res, nil
}
func parseFieldForJson(field compare.Field, status string) (any, error) {
	var value any

	if status == compare.Added {
		value = field.NewValue
	} else if status == compare.Removed {
		value = field.OldValue
	}

	if value == nil && field.Children != nil {
		value = field.Children
	}

	if _, ok := value.([]compare.Field); ok {
		childrenData := map[string]any{}
		for _, v := range value.([]compare.Field) {
			childData, err := parseFieldForJson(v, v.Status)
			if err != nil {
				return nil, err
			}
			childrenData[v.Name] = childData
		}

		children := map[string]any{}
		for k, v := range childrenData {
			children[k] = v
		}

		return children, nil
	} else {
		return value, nil
	}
}
