package utils

import (
	"os"

	"encoding/json"

	"gopkg.in/yaml.v3"
)

func DecodeYAML(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, dst); err != nil {
		return err
	}

	return nil
}

func DecodeJSON(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return err
	}

	return nil
}
