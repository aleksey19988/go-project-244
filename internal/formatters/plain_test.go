package formatters

import (
	"code/internal/diff"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainFormatter_Format(t *testing.T) {
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

	fields = []diff.Field{
		{
			Name:     "Name 1",
			Status:   Removed,
			OldValue: 123,
			NewValue: "new value",
			Depth:    0,
			Children: nil,
		},
	}
	_, err = pf.Format(fields)
	assert.Error(t, err)
	assert.Equal(t, "deep is negative or zero", err.Error())
}

func TestFormatValue(t *testing.T) {
	assert.Equal(t, "null", FormatValue(nil))
	assert.Equal(t, "'str'", FormatValue("str"))
	assert.Equal(t, "false", FormatValue(false))
	assert.Equal(t, "[complex value]", FormatValue([]diff.Field{}))
	assert.Equal(t, "WARNING: unknown type: int\n", FormatValue(123))
}
func TestPlainFormatter_GetOutputString(t *testing.T) {
	t.Run("Zero deep error", func(t *testing.T) {
		fld := diff.Field{}
		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		_, err = f.GetOutputString(fld)
		assert.Error(t, err)
		assert.Equal(t, "deep is negative or zero", err.Error())

		fld = diff.Field{
			Name:     "Name 2",
			Status:   Added,
			OldValue: nil,
			NewValue: "456",
			Depth:    0,
			Children: []diff.Field{
				{
					Name:  "Name 1",
					Depth: 0,
				},
			},
		}

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

		fld = diff.Field{
			Name:     "Name 3",
			Status:   Added,
			OldValue: nil,
			NewValue: nil,
			Depth:    1,
			Children: []diff.Field{},
		}
		res, err = f.GetOutputString(fld)
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was added with value: %v\n", fld.Name, "[complex value]")
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
	t.Run("Update field", func(t *testing.T) {
		fld := diff.Field{
			Name:     "Name 3",
			Status:   Updated,
			OldValue: "123",
			NewValue: "456",
			Depth:    1,
			Children: nil,
		}

		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was updated. From '%v' to '%v'\n", fld.Name, fld.OldValue, fld.NewValue)
		assert.Equal(t, expected, res)

		fld = diff.Field{
			Name:     "Name 3",
			Status:   Updated,
			OldValue: "123",
			NewValue: []diff.Field{},
			Depth:    1,
			Children: nil,
		}

		res, err = f.GetOutputString(fld)
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From '%v' to %v\n", fld.Name, fld.OldValue, "[complex value]")
		assert.Equal(t, expected, res)

		fld = diff.Field{
			Name:     "Name 3",
			Status:   Updated,
			OldValue: []diff.Field{},
			NewValue: "123",
			Depth:    1,
			Children: nil,
		}

		res, err = f.GetOutputString(fld)
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From %v to '%v'\n", fld.Name, "[complex value]", fld.NewValue)
		assert.Equal(t, expected, res)

		fld = diff.Field{
			Name:     "Name 3",
			Status:   Updated,
			OldValue: nil,
			NewValue: "123",
			Depth:    1,
			Children: []diff.Field{},
		}

		res, err = f.GetOutputString(fld)
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From %v to '%v'\n", fld.Name, "[complex value]", fld.NewValue)
		assert.Equal(t, expected, res)

		fld = diff.Field{
			Name:     "Name 3",
			Status:   Updated,
			OldValue: "456",
			NewValue: nil,
			Depth:    1,
			Children: []diff.Field{},
		}

		res, err = f.GetOutputString(fld)
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From '%v' to %v\n", fld.Name, fld.OldValue, "[complex value]")
		assert.Equal(t, expected, res)
	})
	t.Run("Deep is 2", func(t *testing.T) {
		fld := diff.Field{
			Name:     "parent_name",
			Status:   "",
			OldValue: nil,
			NewValue: "456",
			Depth:    1,
			Children: []diff.Field{
				{
					Name:     "child_name",
					Status:   Added,
					OldValue: nil,
					NewValue: "new value",
					Depth:    2,
					Children: nil,
				},
			},
		}

		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was added with value: '%v'\n", fmt.Sprintf("%s.%s", fld.Name, fld.Children[0].Name), fld.Children[0].NewValue)
		assert.Equal(t, expected, res)
	})
}
