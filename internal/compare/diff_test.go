package compare

import (
	fls "code/internal/files"
	"code/internal/parser"
	"code/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	_, err := storage.GetFilesData("", "")
	require.Error(t, err)
	assert.Equal(t, "failed to read : path cannot be empty", err.Error())

	_, err = storage.GetFilesData("path.json", "")
	require.Error(t, err)
	assert.Equal(t, "failed to read path.json: open path.json: no such file or directory", err.Error())

	_, err = storage.GetFilesData("/path/to/file", "")
	require.Error(t, err)
	assert.Equal(t, "failed to read /path/to/file: open /path/to/file: no such file or directory", err.Error())

	_, err = storage.GetFilesData("/path/to/file.json", "/path/to/file.yaml")
	require.Error(t, err)
	assert.Equal(t, "failed to read /path/to/file.json: open /path/to/file.json: no such file or directory", err.Error())
}
func TestDiff(t *testing.T) {
	t.Run("Test json", func(t *testing.T) {
		files, err := storage.GetFilesData("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
		require.NoError(t, err)

		err = fls.Validate(files)
		if err != nil {
			require.NoError(t, err)
		}

		var filesData []map[string]any
		for _, f := range files {
			fileData, err := parser.ParseData(*f)
			require.NoError(t, err)
			filesData = append(filesData, fileData)
		}

		fields := GenDiff(filesData[0], filesData[1])
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
		files, err := storage.GetFilesData("../../testdata/fixture/yaml/file1.yaml", "../../testdata/fixture/yaml/file2.yaml")
		require.NoError(t, err)

		err = fls.Validate(files)
		if err != nil {
			require.NoError(t, err)
		}

		var filesData []map[string]any
		for _, f := range files {
			fileData, err := parser.ParseData(*f)
			require.NoError(t, err)
			filesData = append(filesData, fileData)
		}

		fields := GenDiff(filesData[0], filesData[1])
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
