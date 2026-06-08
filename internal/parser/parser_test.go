package parser

import (
	"code/internal/formatters"
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
	f := formatters.Field{
		Name:         "field-name",
		TypeOfChange: Added,
		Value:        "field-value",
		Deep:         1,
		Children:     nil,
	}
	formatter, err := formatters.NewFormatter(formatters.DefaultOutputFormat)
	require.NoError(t, err)

	res, err := formatter.GetOutputString(f)
	require.NoError(t, err)
	assert.Equal(t, "  + field-name: field-value\n", res)
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

		expectedField := formatters.Field{
			Name:         "group2",
			TypeOfChange: "-",
			Deep:         1,
			Children: []formatters.Field{
				formatters.Field{
					Name:  "abc",
					Value: float64(12345),
					Deep:  2,
				},
				formatters.Field{
					Name: "deep",
					Deep: 2,
					Children: []formatters.Field{
						formatters.Field{
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
		expectedField := formatters.Field{
			Name:         "group2",
			TypeOfChange: "-",
			Deep:         1,
			Children: []formatters.Field{
				formatters.Field{
					Name:  "abc",
					Value: 12345,
					Deep:  2,
				},
				formatters.Field{
					Name: "deep",
					Deep: 2,
					Children: []formatters.Field{
						formatters.Field{
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
