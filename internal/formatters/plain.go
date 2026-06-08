package formatters

import (
	"errors"
	"fmt"
)

type PlainFormatter struct{}

func (pf *PlainFormatter) GetOutput(fields []Field) (string, error) {
	res := ""

	for _, field := range fields {
		fieldData, err := pf.GetOutputString(field)
		if err != nil {
			return "", err
		}
		res += fieldData
	}

	return res, nil
}
func (pf *PlainFormatter) GetOutputString(f Field) (string, error) {
	res := ""

	if f.Deep <= 0 {
		return "", errors.New("deep is negative or zero")
	}

	if f.Children != nil {
		for _, child := range f.Children {
			child.Name = fmt.Sprintf("%s.%s", f.Name, child.Name)
			childData, err := pf.GetOutputString(child)
			if err != nil {
				return "", err
			}
			res += childData
		}
	} else {

		switch f.TypeOfChange {
		case Added:
			if f.Value == nil {
				f.Value = "null"
			}

			res = fmt.Sprintf("Property '%s' was added with value: %v\n", f.Name, f.Value)
		case Removed:
			res = fmt.Sprintf("Property '%s' was removed\n", f.Name)
		default:
			if f.Value == nil {
				f.Value = "null"
			}

			res = fmt.Sprintf("Property '%s' was updated. From %v to %v\n", f.Name, f.Value, "...")
		}
	}

	return res, nil
}
