package formatters

import (
	"code/internal/diff"
	"errors"
)

const (
	Added               = "+"
	Removed             = "-"
	Updated             = "~"
	DefaultOutputFormat = "stylish"
)

type Formatter interface {
	Format(fields []diff.Field) (string, error)
	GetOutputString(f diff.Field) (string, error)
}

func NewFormatter(format string) (Formatter, error) {
	if format == "" {
		format = DefaultOutputFormat
	}
	switch format {
	case "stylish":
		return &StylishFormatter{}, nil
	case "plain":
		return &PlainFormatter{}, nil
	}
	return nil, errors.New("format is not supported")
}
