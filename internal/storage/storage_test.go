package storage

import (
	fls "code/internal/files"
	"code/internal/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_LoadFiles(t *testing.T) {
	_, err := GetFilesData("path1.json", "path2.json")
	require.Error(t, err)

	_, err = GetFilesData("", "path2.json")
	require.Error(t, err)
	assert.Equal(t, "failed to read : path cannot be empty", err.Error())

	_, err = GetFilesData("path1.json", "path2.yaml")
	require.Error(t, err)
	assert.Equal(t, "failed to read path1.json: open path1.json: no such file or directory", err.Error())

	files, err := GetFilesData("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
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
	assert.NoError(t, err)
	assert.Equal(t, 2, len(filesData))
}

func TestStorage_CreateMapsFromData(t *testing.T) {
	t.Run("Invalid file to create maps", func(t *testing.T) {
		_, err := GetFilesData("../../testdata/fixture/file.txt", "../../testdata/fixture/json/file2.json")
		require.NoError(t, err)
	})
	t.Run("Valid json files", func(t *testing.T) {
		files, err := GetFilesData("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
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
		assert.Equal(t, 2, len(filesData))

		m1 := filesData[0]
		expected := map[string]any{
			"abc": float64(12345),
			"deep": map[string]any{
				"id": float64(45),
			},
		}
		assert.Equal(t, expected, m1["group2"])
		m2 := filesData[1]
		expected = map[string]any{
			"foo":  "bar",
			"baz":  "bars",
			"nest": "str",
		}
		assert.Equal(t, expected, m2["group1"])
	})
	t.Run("Valid yaml files", func(t *testing.T) {
		files, err := GetFilesData("../../testdata/fixture/yaml/file1.yaml", "../../testdata/fixture/yaml/file2.yaml")
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
		assert.Equal(t, 2, len(filesData))

		m1 := filesData[0]
		expected := map[string]any{
			"abc": 12345,
			"deep": map[string]any{
				"id": 45,
			},
		}
		assert.Equal(t, expected, m1["group2"])
		m2 := filesData[1]
		expected = map[string]any{
			"foo":  "bar",
			"baz":  "bars",
			"nest": "str",
		}
		assert.Equal(t, expected, m2["group1"])
	})
}
