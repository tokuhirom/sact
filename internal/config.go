package internal

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DefaultZone string `toml:"default_zone"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "sact", "config.toml")

	config := &Config{
		DefaultZone: "tk1b",
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return nil, err
	}

	return config, nil
}
