package parser

import (
	"code/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	err := Validate("")
	require.Error(t, err)
	assert.Equal(t, "path to file is empty", err.Error())

	err = Validate("/path/to/file")
	require.Error(t, err)
	assert.Equal(t, "path to file must have .json extension", err.Error())

	err = Validate("/path/to/file.json")
	require.NoError(t, err)
}
func TestGetOutputString(t *testing.T) {
	assert.Equal(t, "  - name: value\n", GetOutputString("-", "name", "value"))
}
func TestFormatOutput(t *testing.T) {
	t.Run("Test add value", func(t *testing.T) {
		f := Field{
			Name:         "Test name",
			TypeOfChange: "+",
			OldValue:     nil,
			NewValue:     195,
		}

		output := FormatOutput([]Field{f}, "")
		assert.Equal(t, "{\n  + Test name: 195\n}", output)
	})

	t.Run("Test remove value", func(t *testing.T) {
		f := Field{
			Name:         "Some name",
			TypeOfChange: "-",
			OldValue:     "some value",
			NewValue:     nil,
		}

		output := FormatOutput([]Field{f}, "")
		assert.Equal(t, "{\n  - Some name: some value\n}", output)
	})

	t.Run("Equal values", func(t *testing.T) {
		f := Field{
			Name:         "Title",
			TypeOfChange: "",
			OldValue:     false,
			NewValue:     false,
		}

		output := FormatOutput([]Field{f}, "")
		assert.Equal(t, "{\n    Title: false\n}", output)
	})
}
func TestDiff(t *testing.T) {
	s := storage.NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
	err := s.SetRawData()
	require.NoError(t, err)

	maps, err := s.CreateMapsFromData()
	require.NoError(t, err)

	fields, err := Diff(maps)
	require.NoError(t, err)
	assert.Equal(t, 6, len(fields))
	assert.Equal(t, Field{
		Name:         "host",
		TypeOfChange: "",
		OldValue:     "hexlet.io",
		NewValue:     "hexlet.io",
	}, fields[1])
}
func TestParse(t *testing.T) {
	s := storage.NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
	res, err := Parse(s, "")
	require.NoError(t, err)
	assert.Equal(t, res, "{\n  - follow: false\n    host: hexlet.io\n  - proxy: 123.234.53.22\n  - timeout: 50\n  + timeout: 20\n  + verbose: true\n}")
}
