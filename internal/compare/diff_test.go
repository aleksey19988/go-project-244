package compare

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
}
func TestDiff(t *testing.T) {
	t.Run("Test json", func(t *testing.T) {
		s, err := storage.NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
		require.NoError(t, err)

		var filesData []map[string]any
		for _, f := range s.Files {
			fileData, err := f.CreateMapFromData()
			require.NoError(t, err)
			filesData = append(filesData, fileData)
		}

		fields, err := GenDiff(filesData, 1)
		require.NoError(t, err)
		assert.Equal(t, 4, len(fields))

		expectedField := Field{
			Name:   "group2",
			Status: Removed,
			Depth:  1,
			Children: []Field{
				{
					Name:     "abc",
					OldValue: float64(12345),
					NewValue: float64(12345),
					Depth:    2,
				},
				{
					Name:  "deep",
					Depth: 2,
					Children: []Field{
						{
							Name:     "id",
							OldValue: float64(45),
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

		var filesData []map[string]any
		for _, f := range s.Files {
			fileData, err := f.CreateMapFromData()
			require.NoError(t, err)
			filesData = append(filesData, fileData)
		}

		fields, err := GenDiff(filesData, 1)
		require.NoError(t, err)
		assert.Equal(t, 4, len(fields))
		expectedField := Field{
			Name:   "group2",
			Status: Removed,
			Depth:  1,
			Children: []Field{
				{
					Name:     "abc",
					OldValue: 12345,
					NewValue: 12345,
					Depth:    2,
				},
				{
					Name:  "deep",
					Depth: 2,
					Children: []Field{
						{
							Name:     "id",
							OldValue: 45,
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
