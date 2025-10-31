// Package config provides configuration management for LazyLinux.
// It handles loading, saving, and validating configuration settings stored
// in YAML format at ~/.config/lazylinux/config.yaml.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the LazyLinux configuration
type Config struct {
	PackageManager string `yaml:"package_manager"` // "dnf", "apt", or "pacman"
	FlatpakEnabled bool   `yaml:"flatpak_enabled"` // FlatpakEnabled indicates if Flatpak support is available
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "lazylinux", "config.yaml")
}

// ConfigExists checks if the config file exists
func ConfigExists() bool {
	path := GetConfigPath()
	_, err := os.Stat(path)
	return err == nil
}

// LoadConfig reads the config file
func LoadConfig() (*Config, error) {
	path := GetConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}

	return &cfg, nil
}

// SaveConfig writes the config file
func SaveConfig(cfg *Config) error {
	path := GetConfigPath()

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("could not create config directory: %v", err)
	}

	// Convert config to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not convert config to YAML: %v", err)
	}

	// Write to file
	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return fmt.Errorf("could not write config file: %v", err)
	}

	return nil
}
