package parser

import (
	"errors"
	"os"
	"path/filepath"
)

func Parse(path1, path2, format string) (string, error) {
	if path1 == "" {
		return "", errors.New("Path to file 1 is empty")
	}

	if path2 == "" {
		return "", errors.New("Path to file 2 is empty")
	}

	ext1 := filepath.Ext(path1)
	if ext1 != "json" {
		return "", errors.New("Path to file 1 must have .json extension")
	}

	ext2 := filepath.Ext(path2)
	if ext2 != "json" {
		return "", errors.New("Path to file 2 must have .json extension")
	}

	data1, err := os.ReadFile(path1)
	if err != nil {
		return "", err
	}
	return string(data1), nil
}
