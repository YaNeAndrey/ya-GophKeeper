package config

import (
	"encoding/json"
	"os"
)

type Config struct {
}

func ParseConfigFromJSON(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := Config{}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
