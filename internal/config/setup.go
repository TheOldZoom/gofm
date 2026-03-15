package config

import (
	"bufio"
	"errors"
	"fmt"
	"gofm/internal/api"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

func RunSetup() (*Config, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("First, let's get your Last.fm username.")
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')

	username = strings.TrimSpace(username)
	fmt.Println("Now let's get your Last.fm API key. Check https://github.com/theOldZoom/gofm#api-key for more information.")
	fmt.Print("API Key: ")
	apikey, _ := reader.ReadString('\n')
	apikey = strings.TrimSpace(apikey)

	cfg := Config{
		Username: username,
		ApiKey:   apikey,
	}

	if err := ValidateAPIKey(cfg.ApiKey); err != nil {
		fmt.Printf("Invalid API key: %s\n", validationMessage(err))
		return RunSetup()
	}

	if err := ValidateUsername(cfg.Username, cfg.ApiKey); err != nil {
		fmt.Printf("Invalid username: %s\n", validationMessage(err))
		return RunSetup()
	}

	if err := Save(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := Path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func ValidateAPIKey(apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("API key is required")
	}

	return api.ValidateAPIKey(apiKey)
}

func ValidateUsername(username string, apiKey string) error {
	if strings.TrimSpace(username) == "" {
		return fmt.Errorf("username is required")
	}

	return api.ValidateUsername(username, apiKey)
}

func validationMessage(err error) string {
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Message
	}

	return err.Error()
}
