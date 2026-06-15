package formatters

import (
	"code/internal/diff"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStylishFormatter_Format(t *testing.T) {
	fields := []diff.Field{
		{
			Name:     "field_1",
			Status:   Added,
			OldValue: nil,
			NewValue: 123,
			Depth:    1,
			Children: nil,
		},
		{
			Name:     "field_2",
			Status:   Removed,
			OldValue: 456,
			NewValue: nil,
			Depth:    1,
			Children: nil,
		},
	}
	formatter, err := NewFormatter("")
	assert.NoError(t, err)
	assert.NotNil(t, formatter)

	res, err := formatter.Format(fields)
	assert.NoError(t, err)
	expected := fmt.Sprintf("{\n  + %s: %v\n  - %s: %v\n}", fields[0].Name, fields[0].NewValue, fields[1].Name, fields[1].OldValue)
	assert.Equal(t, expected, res)
}

func TestStylishFormatter_GetOutputString(t *testing.T) {
	t.Run("Zero deep error", func(t *testing.T) {
		fld := diff.Field{}
		f, err := NewFormatter("")
		assert.NoError(t, err)
		assert.NotNil(t, f)

		_, err = f.GetOutputString(fld)
		assert.Error(t, err)
		assert.Equal(t, "deep is negative or zero", err.Error())
	})
	t.Run("Add field", func(t *testing.T) {
		fld := diff.Field{
			Name:     "Name 1",
			Status:   Added,
			OldValue: nil,
			NewValue: 123,
			Depth:    1,
			Children: nil,
		}
		f, err := NewFormatter("")
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("  + %s: %v\n", fld.Name, fld.NewValue)
		assert.Equal(t, expected, res)

		fld = diff.Field{
			Name:     "Name 1",
			Status:   Added,
			OldValue: 123,
			NewValue: nil,
			Depth:    1,
			Children: nil,
		}

		res, err = f.GetOutputString(fld)
		assert.NoError(t, err)
		expected = fmt.Sprintf("  + %s: %v\n", fld.Name, "null")
		assert.Equal(t, expected, res)
	})
	t.Run("Add field with children", func(t *testing.T) {
		fld := diff.Field{
			Name:     "Name 1",
			Status:   "",
			OldValue: 123,
			NewValue: nil,
			Depth:    1,
			Children: []diff.Field{
				diff.Field{
					Name:     "child_name",
					Status:   Added,
					OldValue: nil,
					NewValue: "888",
					Depth:    2,
					Children: nil,
				},
			},
		}

		f, err := NewFormatter("")
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("    %s: {\n      + %s: %v\n    }\n", fld.Name, fld.Children[0].Name, fld.Children[0].NewValue)
		assert.Equal(t, expected, res)
	})
	t.Run("Update field", func(t *testing.T) {
		fld := diff.Field{
			Name:     "Name 1",
			Status:   Updated,
			OldValue: 123,
			NewValue: 456,
			Depth:    1,
			Children: nil,
		}
		f, err := NewFormatter("")
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("  - %s: %v\n  + %s: %v\n", fld.Name, fld.OldValue, fld.Name, fld.NewValue)
		assert.Equal(t, expected, res)
	})
	t.Run("Update field with children", func(t *testing.T) {
		fld := diff.Field{
			Name:     "parent_name",
			Status:   Updated,
			OldValue: 123,
			NewValue: []diff.Field{
				diff.Field{
					Name:     "new_field_name",
					Status:   "",
					OldValue: 444,
					NewValue: 555,
					Depth:    2,
					Children: nil,
				},
			},
			Depth: 1,
			Children: []diff.Field{
				diff.Field{
					Name:     "child_1",
					Status:   "",
					OldValue: 777,
					NewValue: nil,
					Depth:    2,
					Children: nil,
				},
			},
		}

		f, err := NewFormatter("")
		assert.NoError(t, err)
		assert.NotNil(t, f)

		res, err := f.GetOutputString(fld)
		assert.NoError(t, err)
		expected := fmt.Sprintf("  - %s: {\n        %s: %v\n    }\n  + %s: {\n        %s: %v\n    }\n", fld.Name, fld.Children[0].Name, fld.Children[0].OldValue, fld.Name, fld.NewValue.([]diff.Field)[0].Name, fld.NewValue.([]diff.Field)[0].OldValue)
		assert.Equal(t, expected, res)

		fld = diff.Field{
			Name:     "parent_name",
			Status:   Updated,
			OldValue: 123,
			NewValue: 45678,
			Depth:    1,
			Children: []diff.Field{
				diff.Field{
					Name:     "child_1",
					Status:   "",
					OldValue: 777,
					NewValue: nil,
					Depth:    2,
					Children: nil,
				},
			},
		}
		res, err = f.GetOutputString(fld)
		assert.NoError(t, err)
		expected = fmt.Sprintf("  - %s: {\n        %s: %v\n    }\n  + %s: %v\n", fld.Name, fld.Children[0].Name, fld.Children[0].OldValue, fld.Name, fld.NewValue)
		assert.Equal(t, expected, res)
	})
}
