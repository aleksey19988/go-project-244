package formatters

import (
	"errors"
)

const (
	Added               = "+"
	Removed             = "-"
	DefaultOutputFormat = "stylish"
)

type Field struct {
	Name         string
	TypeOfChange string
	Value        any
	Deep         int
	Children     []Field
}

type Formatter interface {
	GetOutput(fields []Field) (string, error)
	GetOutputString(f Field) (string, error)
}

func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "stylish":
		return &StylishFormatter{}, nil
	case "plain":
		return &PlainFormatter{}, nil
	}
	return nil, errors.New("format is not supported")
}
