package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/epheo/anytype-go/pkg/anytype"
)

// Configuration constants
const (
	configFileName = "anytype_auth.json"
	configDirName  = "anytype-go"
	configFileMode = 0600
	configDirMode  = 0755
)

var (
	ErrConfigNotFound = errors.New("configuration file not found")
	ErrInvalidConfig  = errors.New("invalid configuration")
)

// getConfigFilePath returns the path to the config file
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	// Create config directory path
	configDir := filepath.Join(homeDir, ".config", configDirName)

	// Ensure the directory exists
	if err := os.MkdirAll(configDir, configDirMode); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}

	return filepath.Join(configDir, configFileName), nil
}

// loadAuthConfig loads authentication config from file
func loadAuthConfig() (*anytype.AuthConfig, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the JSON
	var config anytype.AuthConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// saveAuthConfig saves authentication config to file
func saveAuthConfig(config *anytype.AuthConfig) error {
	if config == nil {
		return ErrInvalidConfig
	}

	// Validate config before saving
	if err := validateConfig(config); err != nil {
		return err
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Convert to JSON with proper indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(configPath, data, configFileMode); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// validateConfig checks if the configuration is valid
func validateConfig(config *anytype.AuthConfig) error {
	if config == nil {
		return ErrInvalidConfig
	}

	if config.ApiURL == "" {
		return fmt.Errorf("%w: missing API URL", ErrInvalidConfig)
	}

	if config.SessionToken == "" {
		return fmt.Errorf("%w: missing session token", ErrInvalidConfig)
	}

	if config.AppKey == "" {
		return fmt.Errorf("%w: missing app key", ErrInvalidConfig)
	}

	if config.Timestamp.IsZero() {
		return fmt.Errorf("%w: missing timestamp", ErrInvalidConfig)
	}

	return nil
}

// RemoveConfig deletes the configuration file
func RemoveConfig() error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.Remove(configPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing config file: %w", err)
	}

	return nil
}
