package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_GetPaths(t *testing.T) {
	s := NewStorage("path1", "path2")
	assert.Equal(t, 2, len(s.GetPaths()))
	assert.Equal(t, "path1", s.GetPaths()[0])
	assert.Equal(t, "path2", s.GetPaths()[1])
}

func TestStorage_SetRawData(t *testing.T) {
	s := NewStorage("path1", "path2")
	err := s.SetRawData()
	assert.Error(t, err)

	s = NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
	err = s.SetRawData()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(s.RawData))
}

func TestStorage_CreateMapsFromData(t *testing.T) {
	t.Run("Invalid file to create maps", func(t *testing.T) {
		s := NewStorage("../../testdata/fixture/file.txt", "../../testdata/fixture/json/file2.json")
		err := s.SetRawData()
		assert.NoError(t, err)
		_, err = s.CreateMapsFromData()
		assert.Error(t, err)
	})

	t.Run("Valid files", func(t *testing.T) {
		s := NewStorage("../../testdata/fixture/json/file1.json", "../../testdata/fixture/json/file2.json")
		err := s.SetRawData()
		assert.NoError(t, err)
		maps, err := s.CreateMapsFromData()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(maps))

		m1 := maps[0]
		assert.Equal(t, false, m1["follow"])
		assert.Equal(t, "123.234.53.22", m1["proxy"])
		assert.Equal(t, float64(50), m1["timeout"])

		m2 := maps[1]
		assert.Equal(t, true, m2["verbose"])
		assert.Equal(t, float64(20), m2["timeout"])
		assert.Equal(t, "hexlet.io", m2["host"])
	})
}
