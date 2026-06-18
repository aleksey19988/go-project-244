package formatters

import (
	"code/internal/compare"
	"errors"
)

const (
	DefaultOutputFormat = StylishOutputFormat
	StylishOutputFormat = "stylish"
	PlainOutputFormat   = "plain"
	JsonOutputFormat    = "json"
)

type Formatter interface {
	Format(fields []compare.Field) (string, error)
}

func NewFormatter(format string) (Formatter, error) {
	if format == "" {
		format = DefaultOutputFormat
	}
	switch format {
	case StylishOutputFormat:
		return new(StylishFormatter), nil
	case PlainOutputFormat:
		return new(PlainFormatter), nil
	case JsonOutputFormat:
		return new(JsonFormatter), nil
	}
	return nil, errors.New("format is not supported")
}
