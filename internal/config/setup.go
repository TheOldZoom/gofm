package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theOldZoom/gofm/internal/api"
	"github.com/theOldZoom/gofm/internal/verbose"

	"go.yaml.in/yaml/v3"
)

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

	err = os.WriteFile(path, data, 0644)
	if err == nil {
		verbose.Printf("saved config file: %s", path)
	}
	return err
}

func ValidateAPIKey(apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("API key is required")
	}

	verbose.Printf("validating api key")
	return api.ValidateAPIKey(apiKey)
}

func ValidateUsername(username string, apiKey string) error {
	if strings.TrimSpace(username) == "" {
		return fmt.Errorf("username is required")
	}

	verbose.Printf("validating username: %s", username)
	return api.ValidateUsername(username, apiKey)
}

func ValidationMessage(err error) string {
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Message
	}
	return err.Error()
}
