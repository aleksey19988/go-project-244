package formatters

import (
	"code/internal/diff"
	"encoding/json"
	"fmt"
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

func (jf JsonFormatter) getFieldDataForOutput(fields []diff.Field) (map[string]map[string]any, error) {
	res := make(map[string]map[string]any)

	for _, field := range fields {
		fmt.Println(field)
		fieldData := map[string]any{}

		if field.Status != "" {
			fieldData["status"] = field.Status
		}

		switch field.Status {
		case Updated:
			if _, ok := field.OldValue.([]diff.Field); ok {
				children, err := jf.getFieldDataForOutput(field.OldValue.([]diff.Field))
				if err != nil {
					return nil, err
				}
				fieldData["old_value"] = children
			} else {
				fieldData["old_value"] = field.OldValue
			}

			if _, ok := field.NewValue.([]diff.Field); ok {
				children, err := jf.getFieldDataForOutput(field.NewValue.([]diff.Field))
				if err != nil {
					return nil, err
				}
				fieldData["new_value"] = children
			} else {
				fieldData["new_value"] = field.NewValue
			}
		case Added:
			if field.NewValue != nil {
				fieldData["value"] = field.NewValue
			} else if _, ok := field.NewValue.([]diff.Field); ok {
				childrenData, err := jf.getFieldDataForOutput(field.Children)
				if err != nil {
					fmt.Println(err.Error())
					return map[string]map[string]any{}, err
				}
				for fieldName, v := range childrenData {
					fieldData[fieldName] = v
				}
			}
		case Removed:
			fieldData["value"] = field.OldValue
		case Unchanged:
			fieldData["value"] = field.OldValue

			res[field.Name] = fieldData
		default:
			if field.Children != nil {

			}
		}
	}

	return res, nil
}
