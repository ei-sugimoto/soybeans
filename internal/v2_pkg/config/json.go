package config

import (
	"fmt"
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func Load(path string) (*TConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config TConfig
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &config, nil

}
