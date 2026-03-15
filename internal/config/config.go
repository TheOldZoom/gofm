package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Username string `yaml:"username"`
	ApiKey   string `yaml:"api_key"`
}

var configPath string

func SetPath(path string) {
	configPath = path
}

func Path() (string, error) {
	if configPath != "" {
		return configPath, nil
	}

	if path := viper.ConfigFileUsed(); path != "" {
		return path, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "gofm", "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)

	return &config, err
}
