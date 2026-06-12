package formatters

import (
	"code/internal/diff"
	"fmt"
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
	expected := fmt.Sprintf("Property '%s' was removed\n", fields[0].Name)
	assert.Equal(t, expected, res)
}

func TestGetOutputString(t *testing.T) {
	t.Run("Zero deep error", func(t *testing.T) {
		fld := diff.Field{
			Name:     "Name 2",
			Status:   Added,
			OldValue: nil,
			NewValue: "456",
			Depth:    0,
			Children: nil,
		}
		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		_, err = f.GetOutputString(fld)
		assert.Error(t, err)
		assert.Equal(t, "deep is negative or zero", err.Error())
	})

	t.Run("Add field", func(t *testing.T) {
		fld := diff.Field{
			Name:     "Name 3",
			Status:   Added,
			OldValue: nil,
			NewValue: "456",
			Depth:    1,
			Children: nil,
		}

		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was added with value: '%v'\n", fld.Name, fld.NewValue)
		assert.Equal(t, expected, res)
	})

	t.Run("Remove field", func(t *testing.T) {
		fld := diff.Field{
			Name:     "Name 2",
			Status:   Removed,
			OldValue: nil,
			NewValue: "456",
			Depth:    1,
			Children: nil,
		}

		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was removed\n", fld.Name)
		assert.Equal(t, expected, res)
	})
}
