package formatters

import (
	"code/internal/compare"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainFormatter_Format(t *testing.T) {
	fds := []compare.Field{
		{
			Name:     "Name 1",
			Status:   compare.Removed,
			OldValue: 123,
			NewValue: "new value",
			Depth:    1,
			Children: nil,
		},
	}
	pf, err := NewFormatter(PlainOutputFormat)
	assert.NoError(t, err)
	assert.NotNil(t, pf)

	res, err := pf.Format(fds)
	assert.NoError(t, err)
	expected := fmt.Sprintf("Property '%s' was removed\n", fds[0].Name)
	assert.Equal(t, expected, res)

	fds = []compare.Field{
		{
			Name:     "Name 1",
			Status:   compare.Removed,
			OldValue: 123,
			NewValue: "new value",
			Depth:    0,
			Children: nil,
		},
	}
	_, err = pf.Format(fds)
	assert.Error(t, err)
	assert.Equal(t, "depth is negative or zero", err.Error())
}

func TestFormatValue(t *testing.T) {
	assert.Equal(t, "null", formatValue(nil))
	assert.Equal(t, "'str'", formatValue("str"))
	assert.Equal(t, "false", formatValue(false))
	assert.Equal(t, "[complex value]", formatValue([]compare.Field{}))
	assert.Equal(t, "WARNING: unknown type: int\n", formatValue(123))
}

func TestPlainFormatter_FormatCases(t *testing.T) {
	t.Run("Zero deep error", func(t *testing.T) {
		fld := compare.Field{}
		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		_, err = f.Format([]compare.Field{fld})
		assert.Error(t, err)
		assert.Equal(t, "depth is negative or zero", err.Error())

		fld = compare.Field{
			Name:     "Name 2",
			Status:   compare.Added,
			OldValue: nil,
			NewValue: "456",
			Depth:    0,
			Children: []compare.Field{
				{
					Name:  "Name 1",
					Depth: 0,
				},
			},
		}

		_, err = f.Format([]compare.Field{fld})
		assert.Error(t, err)
		assert.Equal(t, "depth is negative or zero", err.Error())
	})
	t.Run("Add field", func(t *testing.T) {
		fld := compare.Field{
			Name:     "Name 3",
			Status:   compare.Added,
			OldValue: nil,
			NewValue: "456",
			Depth:    1,
			Children: nil,
		}

		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was added with value: '%v'\n", fld.Name, fld.NewValue)
		assert.Equal(t, expected, res)

		fld = compare.Field{
			Name:     "Name 3",
			Status:   compare.Added,
			OldValue: nil,
			NewValue: nil,
			Depth:    1,
			Children: []compare.Field{},
		}
		res, err = f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was added with value: %v\n", fld.Name, "[complex value]")
		assert.Equal(t, expected, res)
	})
	t.Run("Remove field", func(t *testing.T) {
		fld := compare.Field{
			Name:     "Name 2",
			Status:   compare.Removed,
			OldValue: nil,
			NewValue: "456",
			Depth:    1,
			Children: nil,
		}

		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was removed\n", fld.Name)
		assert.Equal(t, expected, res)
	})
	t.Run("Update field", func(t *testing.T) {
		fld := compare.Field{
			Name:     "Name 3",
			Status:   compare.Updated,
			OldValue: "123",
			NewValue: "456",
			Depth:    1,
			Children: nil,
		}

		f, err := NewFormatter(PlainOutputFormat)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was updated. From '%v' to '%v'\n", fld.Name, fld.OldValue, fld.NewValue)
		assert.Equal(t, expected, res)

		fld = compare.Field{
			Name:     "Name 3",
			Status:   compare.Updated,
			OldValue: "123",
			NewValue: []compare.Field{},
			Depth:    1,
			Children: nil,
		}

		res, err = f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From '%v' to %v\n", fld.Name, fld.OldValue, "[complex value]")
		assert.Equal(t, expected, res)

		fld = compare.Field{
			Name:     "Name 3",
			Status:   compare.Updated,
			OldValue: []compare.Field{},
			NewValue: "123",
			Depth:    1,
			Children: nil,
		}

		res, err = f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From %v to '%v'\n", fld.Name, "[complex value]", fld.NewValue)
		assert.Equal(t, expected, res)

		fld = compare.Field{
			Name:     "Name 3",
			Status:   compare.Updated,
			OldValue: nil,
			NewValue: "123",
			Depth:    1,
			Children: []compare.Field{},
		}

		res, err = f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From %v to '%v'\n", fld.Name, "[complex value]", fld.NewValue)
		assert.Equal(t, expected, res)

		fld = compare.Field{
			Name:     "Name 3",
			Status:   compare.Updated,
			OldValue: "456",
			NewValue: nil,
			Depth:    1,
			Children: []compare.Field{},
		}

		res, err = f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected = fmt.Sprintf("Property '%s' was updated. From '%v' to %v\n", fld.Name, fld.OldValue, "[complex value]")
		assert.Equal(t, expected, res)
	})
	t.Run("Deep is 2", func(t *testing.T) {
		fld := compare.Field{
			Name:     "parent_name",
			Status:   "",
			OldValue: nil,
			NewValue: "456",
			Depth:    1,
			Children: []compare.Field{
				{
					Name:     "child_name",
					Status:   compare.Added,
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

		res, err := f.Format([]compare.Field{fld})
		assert.NoError(t, err)
		expected := fmt.Sprintf("Property '%s' was added with value: '%v'\n", fmt.Sprintf("%s.%s", fld.Name, fld.Children[0].Name), fld.Children[0].NewValue)
		assert.Equal(t, expected, res)
	})
}
