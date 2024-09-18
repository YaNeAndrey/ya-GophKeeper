package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	SrvAddr      string
	TempDir      string
	ChunkSize    uint64
	SyncInterval time.Duration
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
