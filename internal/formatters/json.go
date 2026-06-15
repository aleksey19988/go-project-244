package formatters

import "code/internal/diff"

type JsonFormatter struct{}

func (jf JsonFormatter) Format(fields []diff.Field) (string, error) {
	return "", nil
}

func (jf JsonFormatter) GetOutputString(f diff.Field) (string, error) {
	return "", nil
}
