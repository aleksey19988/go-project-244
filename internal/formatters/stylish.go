package formatters

import (
	"code/internal/diff"
	"errors"
	"fmt"
	"strings"
)

type StylishFormatter struct{}

func (sf *StylishFormatter) Format(fields []diff.Field) (string, error) {
	bld := strings.Builder{}
	bld.WriteString("{\n")

	for _, field := range fields {
		fieldData, err := sf.GetOutputString(field)
		if err != nil {
			return "", err
		}
		bld.WriteString(fieldData)
	}

	bld.WriteString("}")

	return bld.String(), nil
}
func (sf *StylishFormatter) GetOutputString(f diff.Field) (string, error) {
	if f.Depth <= 0 {
		return "", errors.New("deep is negative or zero")
	}
	bld := new(strings.Builder)

	marginsCount := f.Depth*4 - 2
	if len(f.Status) == 0 {
		f.Status = " "
	}

	oldValue := f.OldValue
	if oldValue == nil {
		oldValue = "null"
	}

	newValue := f.NewValue
	if newValue == nil {
		newValue = "null"
	}

	if f.Children != nil {
		switch f.Status {
		case Updated:
			_, err := fmt.Fprintf(bld, "%s- %s: {\n", strings.Repeat(" ", marginsCount), f.Name)
			if err != nil {
				return "", err
			}
		case Added:
			_, err := fmt.Fprintf(bld, "%s+ %s: {\n", strings.Repeat(" ", marginsCount), f.Name)
			if err != nil {
				return "", err
			}
		case Removed:
			_, err := fmt.Fprintf(bld, "%s- %s: {\n", strings.Repeat(" ", marginsCount), f.Name)
			if err != nil {
				return "", err
			}
		default:
			_, err := fmt.Fprintf(bld, "%s%s %s: {\n", strings.Repeat(" ", marginsCount), f.Status, f.Name)
			if err != nil {
				return "", err
			}
		}

		for _, child := range f.Children {
			childData, err := sf.GetOutputString(child)
			if err != nil {
				return "", err
			}
			bld.WriteString(childData)
		}

		_, err := fmt.Fprintf(bld, "%s  }\n", strings.Repeat(" ", marginsCount))
		if err != nil {
			return "", err
		}

		if f.Status == Updated {
			if fields, isSlice := newValue.([]diff.Field); isSlice {
				_, err = fmt.Fprintf(bld, "%s+ %s: {\n", strings.Repeat(" ", marginsCount), f.Name)
				if err != nil {
					return "", err
				}
				for _, child := range fields {
					childData, err := sf.GetOutputString(child)
					if err != nil {
						return "", err
					}
					bld.WriteString(childData)
				}
				_, err = fmt.Fprintf(bld, "%s  }\n", strings.Repeat(" ", marginsCount))
				if err != nil {
					return "", err
				}
			} else {
				_, err = fmt.Fprintf(bld, "%s+ %s: %v\n", strings.Repeat(" ", marginsCount), f.Name, newValue)
				if err != nil {
					return "", err
				}
			}
		}
	} else {
		switch f.Status {
		case Updated:
			_, err := fmt.Fprintf(bld, "%s- %s: %v\n", strings.Repeat(" ", marginsCount), f.Name, oldValue)
			if err != nil {
				return "", err
			}
			_, err = fmt.Fprintf(bld, "%s+ %s: %v\n", strings.Repeat(" ", marginsCount), f.Name, newValue)
			if err != nil {
				return "", err
			}
		case Added:
			_, err := fmt.Fprintf(bld, "%s+ %s: %v\n", strings.Repeat(" ", marginsCount), f.Name, newValue)
			if err != nil {
				return "", err
			}
		case Removed:
			_, err := fmt.Fprintf(bld, "%s- %s: %v\n", strings.Repeat(" ", marginsCount), f.Name, oldValue)
			if err != nil {
				return "", err
			}
		default:
			_, err := fmt.Fprintf(bld, "%s%s %s: %v\n", strings.Repeat(" ", marginsCount), f.Status, f.Name, oldValue)
			if err != nil {
				return "", err
			}
		}
	}

	return bld.String(), nil
}
