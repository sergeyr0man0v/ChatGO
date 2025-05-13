package config

import (
	"chatgo/server/internal/db"
	"chatgo/server/internal/services"
	"chatgo/server/router"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Database db.Config       `yaml:"database"`
	Server   router.Config   `yaml:"server"`
	Service  services.Config `yaml:"service"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Could not parse config: %v", err)
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}
