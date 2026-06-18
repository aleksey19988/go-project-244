package storage

import (
	fls "code/internal/files"
	"fmt"
	"os"
	"path/filepath"
)

func GetFilesData(path1, path2 string) (fls.Files, error) {
	var files fls.Files

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

	return files, nil
}
func loadFileData(path string) (*fls.FileData, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("path cannot be empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &fls.FileData{
		RawData:   data,
		Name:      path,
		Extension: filepath.Ext(path),
	}, nil
}
