package formatters

import (
	"code/internal/compare"
	"errors"
	"fmt"
	"strings"
)

type PlainFormatter struct{}

func (pf *PlainFormatter) Format(fields []compare.Field) (string, error) {
	bld := strings.Builder{}

	for _, field := range fields {
		fieldData, err := parseFieldForPlain(field)
		if err != nil {
			return "", err
		}
		bld.WriteString(fieldData)
	}

	return bld.String(), nil
}
func parseFieldForPlain(f compare.Field) (string, error) {
	if f.Depth <= 0 {
		return "", errors.New("depth is negative or zero")
	}

	bld := new(strings.Builder)

	if f.Children != nil && f.Status == "" {
		for _, child := range f.Children {
			cp := child
			cp.Name = fmt.Sprintf("%s.%s", f.Name, child.Name)
			childData, err := parseFieldForPlain(cp)
			if err != nil {
				return "", err
			}
			bld.WriteString(childData)
		}
	} else {
		from := formatValue(f.OldValue)
		to := formatValue(f.NewValue)

		switch f.Status {
		case compare.Added:
			if f.NewValue == nil && f.Children != nil {
				_, err := fmt.Fprintf(bld, "Property '%s' was added with value: %v\n", f.Name, formatValue(f.Children))
				if err != nil {
					return "", err
				}
			} else {
				_, err := fmt.Fprintf(bld, "Property '%s' was added with value: %v\n", f.Name, to)
				if err != nil {
					return "", err
				}
			}
		case compare.Removed:
			_, err := fmt.Fprintf(bld, "Property '%s' was removed\n", f.Name)
			if err != nil {
				return "", err
			}
		case compare.Updated:
			if f.OldValue == nil && f.Children != nil {
				from = formatValue(f.Children)
			}

			if f.NewValue == nil && f.Children != nil {
				to = formatValue(f.Children)
			}

			_, err := fmt.Fprintf(bld, "Property '%s' was updated. From %v to %v\n", f.Name, from, to)
			if err != nil {
				return "", err
			}
		}
	}

	return bld.String(), nil
}
func formatValue(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case []compare.Field, map[string]any:
		return "[complex value]"
	default:
		return fmt.Sprintf("WARNING: unknown type: %T\n", v)
	}
}
