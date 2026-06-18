package parser

import (
	fls "code/internal/files"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func ParseData(fd fls.FileData) (map[string]any, error) {
	result := make(map[string]any)
	switch fd.Extension {
	case fls.JsonExt:
		if err := json.Unmarshal(fd.RawData, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file '%s': %w", fd.Name, err)
		}
	case fls.YamlExt, fls.YmlExt:
		if err := yaml.Unmarshal(fd.RawData, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file '%s': %w", fd.Name, err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension '%s'", fd.Extension)
	}

	return result, nil
}
