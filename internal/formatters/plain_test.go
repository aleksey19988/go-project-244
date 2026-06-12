package formatters

import (
	"code/internal/diff"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	fields := []diff.Field{
		{
			Name:     "Name 1",
			Status:   Removed,
			OldValue: 123,
			NewValue: "new value",
			Depth:    1,
			Children: nil,
		},
	}
	pf, err := NewFormatter(PlainOutputFormat)
	assert.NoError(t, err)
	assert.NotNil(t, pf)

	res, err := pf.Format(fields)
	assert.NoError(t, err)
	assert.Equal(t, "Property 'Name 1' was removed\n", res)
}

func TestGetOutputString(t *testing.T) {
	fld := diff.Field{
		Name:     "Name 2",
		Status:   Added,
		OldValue: nil,
		NewValue: "456",
		Depth:    1,
		Children: nil,
	}

	f, err := NewFormatter(PlainOutputFormat)
	assert.NoError(t, err)
	res, err := f.GetOutputString(fld)
	assert.NoError(t, err)
	assert.Equal(t, "Property 'Name 2' was added with value: '456'\n", res)
}
