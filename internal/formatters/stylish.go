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
		if f.Status == Updated {
			bld.WriteString(fmt.Sprintf("%s%s %s: {\n", strings.Repeat(" ", marginsCount), Removed, f.Name))
		} else {
			bld.WriteString(fmt.Sprintf("%s%s %s: {\n", strings.Repeat(" ", marginsCount), f.Status, f.Name))
		}

		for _, child := range f.Children {
			childData, err := sf.GetOutputString(child)
			if err != nil {
				return "", err
			}
			bld.WriteString(childData)
		}

		bld.WriteString(fmt.Sprintf("%s  }\n", strings.Repeat(" ", marginsCount)))

		if f.Status == Updated {
			if fields, isSlice := newValue.([]diff.Field); isSlice {
				bld.WriteString(fmt.Sprintf("%s%s %s: {\n", strings.Repeat(" ", marginsCount), Added, f.Name))
				for _, child := range fields {
					childData, err := sf.GetOutputString(child)
					if err != nil {
						return "", err
					}
					bld.WriteString(childData)
				}
				bld.WriteString(fmt.Sprintf("%s  }\n", strings.Repeat(" ", marginsCount)))
			} else {
				bld.WriteString(fmt.Sprintf("%s%s %s: %v\n", strings.Repeat(" ", marginsCount), Added, f.Name, newValue))
			}
		}
	} else {
		if f.Status == Updated {
			bld.WriteString(fmt.Sprintf("%s%s %s: %v\n", strings.Repeat(" ", marginsCount), Removed, f.Name, oldValue))
			bld.WriteString(fmt.Sprintf("%s%s %s: %v\n", strings.Repeat(" ", marginsCount), Added, f.Name, newValue))
		} else {
			bld.WriteString(fmt.Sprintf("%s%s %s: %v\n", strings.Repeat(" ", marginsCount), f.Status, f.Name, oldValue))
		}
	}

	return bld.String(), nil
}
