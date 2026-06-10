package parser

import (
	"code/internal/diff"
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
	f := diff.Field{
		Name:     "field-name",
		Status:   formatters.Added,
		NewValue: "field-value",
		Depth:    1,
		Children: nil,
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

		err = s.LoadFiles()
		require.NoError(t, err)

		maps, err := s.CreateMapsFromData()
		require.NoError(t, err)

		fields, err := Diff(maps, 1)
		require.NoError(t, err)
		assert.Equal(t, 4, len(fields))

		expectedField := diff.Field{
			Name:   "group2",
			Status: "-",
			Depth:  1,
			Children: []diff.Field{
				{
					Name:     "abc",
					NewValue: float64(12345),
					Depth:    2,
				},
				{
					Name:  "deep",
					Depth: 2,
					Children: []diff.Field{
						{
							Name:     "id",
							NewValue: float64(45),
							Depth:    3,
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

		err = s.LoadFiles()
		require.NoError(t, err)

		maps, err := s.CreateMapsFromData()
		require.NoError(t, err)

		fields, err := Diff(maps, 1)
		require.NoError(t, err)
		assert.Equal(t, 4, len(fields))
		expectedField := diff.Field{
			Name:   "group2",
			Status: formatters.Removed,
			Depth:  1,
			Children: []diff.Field{
				{
					Name:     "abc",
					NewValue: 12345,
					Depth:    2,
				},
				{
					Name:  "deep",
					Depth: 2,
					Children: []diff.Field{
						{
							Name:     "id",
							NewValue: 45,
							Depth:    3,
						},
					},
				},
			},
		}
		assert.Equal(t, expectedField, fields[2])
	})
}
