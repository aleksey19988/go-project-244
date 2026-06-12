package formatters

import (
	"code/internal/diff"
	"errors"
)

const (
	Added               = "+"
	Removed             = "-"
	Updated             = "~"
	DefaultOutputFormat = StylishOutputFormat
	StylishOutputFormat = "stylish"
	PlainOutputFormat   = "plain"
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
	case StylishOutputFormat:
		return &StylishFormatter{}, nil
	case PlainOutputFormat:
		return &PlainFormatter{}, nil
	}
	return nil, errors.New("format is not supported")
}
