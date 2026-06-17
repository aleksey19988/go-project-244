package formatters

import (
	"code/internal/diff"
	"encoding/json"
)

type JsonFormatter struct{}

func (jf JsonFormatter) Format(fields []diff.Field) (string, error) {
	res, err := jf.getFieldDataForOutput(fields)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (jf JsonFormatter) GetOutputString(f diff.Field) (string, error) {
	return "", nil
}

func (jf JsonFormatter) getFieldDataForOutput(fields []diff.Field) (map[string]any, error) {
	res := make(map[string]any)

	for _, field := range fields {
		fieldData := map[string]any{}

		if field.Status != "" {
			fieldData["status"] = field.Status
		}

		switch field.Status {
		case diff.Added:
			value, err := jf.getParsedData(field, diff.Added)
			if err != nil {
				return nil, err
			}
			fieldData["value"] = value
			res[field.Name] = fieldData
		case diff.Removed, diff.Unchanged:
			value, err := jf.getParsedData(field, diff.Removed)
			if err != nil {
				return nil, err
			}
			fieldData["value"] = value
			res[field.Name] = fieldData
		case diff.Updated:
			oldValue, err := jf.getParsedData(field, diff.Removed)
			if err != nil {
				return nil, err
			}
			fieldData["old_value"] = oldValue

			newValue, err := jf.getParsedData(field, diff.Added)
			if err != nil {
				return nil, err
			}
			fieldData["new_value"] = newValue
			res[field.Name] = fieldData
		default:
			if len(field.Children) > 0 {
				childrenData, err := jf.getFieldDataForOutput(field.Children)
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

func (jf JsonFormatter) getParsedData(field diff.Field, status string) (any, error) {
	var value any

	if status == diff.Added {
		value = field.NewValue
	} else if status == diff.Removed {
		value = field.OldValue
	}

	if value == nil && field.Children != nil {
		value = field.Children
	}

	if _, ok := value.([]diff.Field); ok {
		childrenData, err := jf.getFieldDataForOutput(value.([]diff.Field))
		if err != nil {
			return nil, err
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
