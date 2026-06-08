package formatters

import (
	"errors"
	"fmt"
	"strings"
)

type StylishFormatter struct{}

func (sf *StylishFormatter) GetOutput(fields []Field) (string, error) {
	res := "{\n"

	for _, field := range fields {
		fieldData, err := sf.GetOutputString(field)
		if err != nil {
			return "", err
		}
		res += fieldData
	}

	res += "}"

	return res, nil
}
func (sf *StylishFormatter) GetOutputString(f Field) (string, error) {
	res := ""

	if f.Deep <= 0 {
		return "", errors.New("deep is negative or zero")
	}

	marginsCount := f.Deep*4 - 2
	if len(f.TypeOfChange) == 0 {
		f.TypeOfChange = " "
	}

	if f.Children != nil {
		res = fmt.Sprintf("%s%s %s: {\n", strings.Repeat(" ", marginsCount), f.TypeOfChange, f.Name)

		for _, child := range f.Children {
			childData, err := sf.GetOutputString(child)
			if err != nil {
				return "", err
			}
			res += childData
		}

		res += fmt.Sprintf("%s  }\n", strings.Repeat(" ", marginsCount))
	} else {
		if f.Value == nil {
			f.Value = "null"
		}
		res = fmt.Sprintf("%s%s %s: %v\n", strings.Repeat(" ", marginsCount), f.TypeOfChange, f.Name, f.Value)
	}

	return res, nil
}
