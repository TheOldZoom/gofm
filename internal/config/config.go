package config

import (
	"os"
	"path/filepath"

	"github.com/theOldZoom/gofm/internal/verbose"
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
	verbose.Printf("config path set: %s", path)
}

func Path() (string, error) {
	if configPath != "" {
		verbose.Printf("using cached config path: %s", configPath)
		return configPath, nil
	}

	if path := viper.ConfigFileUsed(); path != "" {
		verbose.Printf("using viper config file path: %s", path)
		return path, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(configDir, "gofm", "config.yaml")
	verbose.Printf("computed default config path: %s", path)
	return path, nil
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
	verbose.Printf("loaded config file: %s", path)
	var config Config
	err = yaml.Unmarshal(data, &config)

	return &config, err
}
