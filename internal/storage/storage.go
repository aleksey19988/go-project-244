package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	YamlExt = ".yaml"
	YmlExt  = ".yml"
	JsonExt = ".json"
)

type Storage struct {
	Files []*FileData
}

type FileData struct {
	Name      string
	Extension string
	RawData   []byte
}

func NewStorage(path1, path2 string) (*Storage, error) {
	extensions := map[string]int{}

	for key, path := range []string{path1, path2} {
		if path == "" {
			return nil, fmt.Errorf("path to file %d is empty", key+1)
		}

		ext := filepath.Ext(path)
		if ext != JsonExt && ext != YamlExt && ext != YmlExt {
			return nil, fmt.Errorf("path to file %d must have .json or .yaml extension", key+1)
		}

		extensions[ext]++
	}

	if len(extensions) > 1 {
		return nil, errors.New("files must have one extension")
	}

	var files []*FileData

	file1, err := loadFileData(path1)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path1, err)
	}
	files = append(files, file1)

	file2, err := loadFileData(path2)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path2, err)
	}
	files = append(files, file2)

	return &Storage{Files: files}, nil
}
func (fd *FileData) CreateMapFromData() (map[string]any, error) {
	result := make(map[string]any)
	switch fd.Extension {
	case JsonExt:
		if err := json.Unmarshal(fd.RawData, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file '%s': %w", fd.Name, err)
		}
	case YamlExt, YmlExt:
		if err := yaml.Unmarshal(fd.RawData, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file '%s': %w", fd.Name, err)
		}
	}

	return result, nil
}
func loadFileData(path string) (*FileData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	return &FileData{
		RawData:   data,
		Name:      path,
		Extension: filepath.Ext(path),
	}, nil
}
