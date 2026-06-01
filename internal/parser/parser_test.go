package parser

import (
	"code/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	_, err := storage.NewStorage("", "")
	require.Error(t, err)
	assert.Equal(t, "path to file 1 is empty", err.Error())

	_, err = storage.NewStorage("path.json", "")
	require.Error(t, err)
	assert.Equal(t, "path to file 2 is empty", err.Error())

	_, err = storage.NewStorage("/path/to/file", "")
	require.Error(t, err)
	assert.Equal(t, "path to file 1 must have .json or .yaml extension", err.Error())

	_, err = storage.NewStorage("/path/to/file.json", "/path/to/file.yaml")
	require.Error(t, err)
	assert.Equal(t, "files must have one extension", err.Error())

	_, err = storage.NewStorage("/path/to/file.json", "/path/to/file.json")
	require.NoError(t, err)

	_, err = storage.NewStorage("/path/to/file.yaml", "/path/to/file.yaml")
	require.NoError(t, err)
}
func TestGetOutputString(t *testing.T) {
	f := Field{
		Name:         "field-name",
		TypeOfChange: Added,
		Value:        "field-value",
		Deep:         1,
		Children:     nil,
	}
	res, err := GetOutputString(f)
	require.NoError(t, err)
	assert.Equal(t, "  + field-name: field-value\n", res)
}
func TestFormatOutput(t *testing.T) {
	t.Run("Test add value", func(t *testing.T) {
		f := Field{
			Name:         "Test name",
			TypeOfChange: "+",
			Value:        195,
		}

		output, err := FormatOutput([]Field{f}, "")
		require.Error(t, err)
		assert.Equal(t, "deep is negative or zero", err.Error())

		f.Deep = 1
		output, err = FormatOutput([]Field{f}, "")
		assert.Equal(t, "{\n  + Test name: 195\n}", output)
	})
	t.Run("Test remove value", func(t *testing.T) {
		f := Field{
			Name:         "Some name",
			TypeOfChange: Removed,
			Value:        "some value",
			Deep:         1,
		}

		output, err := FormatOutput([]Field{f}, "")
		require.NoError(t, err)
		assert.Equal(t, "{\n  - Some name: some value\n}", output)
	})
	t.Run("Equal values", func(t *testing.T) {
		f := Field{
			Name:         "Title",
			TypeOfChange: "",
			Value:        false,
			Deep:         1,
		}

		output, err := FormatOutput([]Field{f}, "")
		require.NoError(t, err)
		assert.Equal(t, "{\n    Title: false\n}", output)
	})
}
func TestDiff(t *testing.T) {
	t.Run("Test json", func(t *testing.T) {
		s, err := storage.NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
		require.NoError(t, err)

		err = s.SetRawData()
		require.NoError(t, err)

		maps, err := s.CreateMapsFromData()
		require.NoError(t, err)

		fields, err := Diff(maps, 1)
		require.NoError(t, err)
		assert.Equal(t, 4, len(fields))

		expectedField := Field{
			Name:         "group2",
			TypeOfChange: "-",
			Deep:         1,
			Children: []Field{
				Field{
					Name:  "abc",
					Value: float64(12345),
					Deep:  2,
				},
				Field{
					Name: "deep",
					Deep: 2,
					Children: []Field{
						Field{
							Name:  "id",
							Value: float64(45),
							Deep:  3,
						},
					},
				},
			},
		}
		assert.Equal(t, expectedField, fields[2])
	})
	t.Run("Test yaml", func(t *testing.T) {
		s, err := storage.NewStorage("../../testdata/fixture/yaml/file1.yaml", "../../testdata/fixture/yaml/file2.yaml")
		require.NoError(t, err)

		err = s.SetRawData()
		require.NoError(t, err)

		maps, err := s.CreateMapsFromData()
		require.NoError(t, err)

		fields, err := Diff(maps, 1)
		require.NoError(t, err)
		assert.Equal(t, 4, len(fields))
		expectedField := Field{
			Name:         "group2",
			TypeOfChange: "-",
			Deep:         1,
			Children: []Field{
				Field{
					Name:  "abc",
					Value: 12345,
					Deep:  2,
				},
				Field{
					Name: "deep",
					Deep: 2,
					Children: []Field{
						Field{
							Name:  "id",
							Value: 45,
							Deep:  3,
						},
					},
				},
			},
		}
		assert.Equal(t, expectedField, fields[2])
	})
}
func TestParse(t *testing.T) {
	s, err := storage.NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
	require.NoError(t, err)

	res, err := Parse(s, "")
	require.NoError(t, err)
	expectedOutput := "{\n    common: {\n      + follow: false\n        setting1: Value 1\n      - setting2: 200\n      - setting3: true\n      + setting3: null\n      + setting4: blah blah\n      + setting5: {\n            key5: value5\n        }\n        setting6: {\n            doge: {\n              - wow: \n              + wow: so much\n            }\n            key: value\n          + ops: vops\n        }\n    }\n    group1: {\n      - baz: bas\n      + baz: bars\n        foo: bar\n      - nest: {\n            key: value\n        }\n      + nest: str\n    }\n  - group2: {\n        abc: 12345\n        deep: {\n            id: 45\n        }\n    }\n  + group3: {\n        deep: {\n            id: {\n                number: 45\n            }\n        }\n        fee: 100500\n    }\n}"
	assert.Equal(t, res, expectedOutput)
}
