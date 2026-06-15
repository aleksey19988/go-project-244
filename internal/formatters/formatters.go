package formatters

import (
	"code/internal/diff"
	"errors"
)

const (
	Added     = "added"
	Removed   = "removed"
	Updated   = "updated"
	Unchanged = "unchanged"

	DefaultOutputFormat = StylishOutputFormat
	StylishOutputFormat = "stylish"
	PlainOutputFormat   = "plain"
	JsonOutputFormat    = "json"
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
		return new(StylishFormatter), nil
	case PlainOutputFormat:
		return new(PlainFormatter), nil
	case JsonOutputFormat:
		return new(JsonFormatter), nil
	}
	return nil, errors.New("format is not supported")
}
