package fls

import (
	"errors"
	"fmt"
)

const (
	YamlExt = ".yaml"
	YmlExt  = ".yml"
	JsonExt = ".json"
)

type Files []*FileData

type FileData struct {
	Name      string
	Extension string
	RawData   []byte
}

func Validate(files Files) error {
	extensions := map[string]int{}

	for key, fd := range files {
		if fd.Name == "" {
			return fmt.Errorf("filename is empty")
		}

		if fd.Extension != JsonExt && fd.Extension != YamlExt && fd.Extension != YmlExt {
			return fmt.Errorf("path to file %d must have .json or .yaml extension", key+1)
		}

		extensions[fd.Extension]++
	}

	if len(extensions) > 1 {
		return errors.New("files must have one extension")
	}

	return nil
}
