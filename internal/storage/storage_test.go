package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_GetPaths(t *testing.T) {
	s, err := NewStorage("path1.json", "path2.json")
	require.NoError(t, err)

	assert.Equal(t, 2, len(s.GetPaths()))
	assert.Equal(t, "path1.json", s.GetPaths()[0])
	assert.Equal(t, "path2.json", s.GetPaths()[1])
}

func TestStorage_SetRawData(t *testing.T) {
	s, err := NewStorage("path1.json", "path2.json")
	require.NoError(t, err)

	err = s.SetRawData()
	assert.Error(t, err)
	_, err = NewStorage("", "path2.json")
	require.Error(t, err)
	assert.Equal(t, "path to file 1 is empty", err.Error())

	_, err = NewStorage("path1.json", "path2.yaml")
	require.Error(t, err)
	assert.Equal(t, "files must have one extension", err.Error())

	s, err = NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
	require.NoError(t, err)
	err = s.SetRawData()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(s.RawData))
}

func TestStorage_CreateMapsFromData(t *testing.T) {
	t.Run("Invalid file to create maps", func(t *testing.T) {
		_, err := NewStorage("../../testdata/fixture/file.txt", "../../testdata/fixture/json/file2.json")
		require.Error(t, err)
		assert.Equal(t, "path to file 1 must have .json or .yaml extension", err.Error())
	})

	t.Run("Valid json files", func(t *testing.T) {
		s, err := NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
		require.NoError(t, err)

		err = s.SetRawData()
		assert.NoError(t, err)
		maps, err := s.CreateMapsFromData()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(maps))

		m1 := maps[0]
		expected := map[string]any{
			"abc": float64(12345),
			"deep": map[string]any{
				"id": float64(45),
			},
		}
		assert.Equal(t, expected, m1["group2"])
		m2 := maps[1]
		expected = map[string]any{
			"foo":  "bar",
			"baz":  "bars",
			"nest": "str",
		}
		assert.Equal(t, expected, m2["group1"])
	})

	t.Run("Valid yaml files", func(t *testing.T) {
		s, err := NewStorage("../../testdata/fixture/yaml/file1.yaml", "../../testdata/fixture/yaml/file2.yaml")
		require.NoError(t, err)

		err = s.SetRawData()
		assert.NoError(t, err)
		maps, err := s.CreateMapsFromData()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(maps))

		m1 := maps[0]
		expected := map[string]any{
			"abc": 12345,
			"deep": map[string]any{
				"id": 45,
			},
		}
		assert.Equal(t, expected, m1["group2"])
		m2 := maps[1]
		expected = map[string]any{
			"foo":  "bar",
			"baz":  "bars",
			"nest": "str",
		}
		assert.Equal(t, expected, m2["group1"])
	})
}
