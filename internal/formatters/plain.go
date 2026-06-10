package formatters

import (
	"code/internal/diff"
	"errors"
	"fmt"
	"strings"
)

type PlainFormatter struct{}

func (pf *PlainFormatter) Format(fields []diff.Field) (string, error) {
	bld := strings.Builder{}

	for _, field := range fields {
		fieldData, err := pf.GetOutputString(field)
		if err != nil {
			return "", err
		}
		bld.WriteString(fieldData)
	}

	return bld.String(), nil
}
func (pf *PlainFormatter) GetOutputString(f diff.Field) (string, error) {
	if f.Depth <= 0 {
		return "", errors.New("deep is negative or zero")
	}

	bld := strings.Builder{}

	if f.Children != nil && f.Status == "" {
		for _, child := range f.Children {
			cp := child
			cp.Name = fmt.Sprintf("%s.%s", f.Name, child.Name)
			childData, err := pf.GetOutputString(cp)
			if err != nil {
				return "", err
			}
			bld.WriteString(childData)
		}
	} else {
		f.OldValue = FormatValue(f.OldValue)
		f.NewValue = FormatValue(f.NewValue)

		switch f.Status {
		case Added:
			bld.WriteString(fmt.Sprintf("Property '%s' was added with value: %v\n", f.Name, f.NewValue))
		case Removed:
			bld.WriteString(fmt.Sprintf("Property '%s' was removed\n", f.Name))
		case Updated:
			bld.WriteString(fmt.Sprintf("Property '%s' was updated. From %v to %v\n", f.Name, f.OldValue, f.NewValue))
		}
	}

	return bld.String(), nil
}

func FormatValue(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case []diff.Field, map[string]any:
		return "[complex value]"
	default:
		return fmt.Sprintf("WARNING: unknown type: %T\n", v)
	}
}
