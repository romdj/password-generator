package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Length           int    `yaml:"length"`
	IncludeUpper     bool   `yaml:"include_upper"`
	IncludeLower     bool   `yaml:"include_lower"`
	IncludeDigits    bool   `yaml:"include_digits"`
	IncludeSymbols   bool   `yaml:"include_symbols"`
	ExcludeAmbiguous bool   `yaml:"exclude_ambiguous"`
	Count            int    `yaml:"count"`
	ShowStrength     bool   `yaml:"show_strength"`
	PolicyTemplate   string `yaml:"policy_template"`
}

func DefaultConfig() Config {
	return Config{
		Length:           12,
		IncludeUpper:     true,
		IncludeLower:     true,
		IncludeDigits:    true,
		IncludeSymbols:   false,
		ExcludeAmbiguous: false,
		Count:            1,
		ShowStrength:     false,
		PolicyTemplate:   "",
	}
}

func LoadConfig() (Config, error) {
	config := DefaultConfig()

	// Load from config files (in order of precedence)
	configPaths := []string{
		".pwgen.yaml",
		".pwgen.yml",
	}

	// Add home directory config paths
	if homeDir, err := os.UserHomeDir(); err == nil {
		configPaths = append(configPaths,
			filepath.Join(homeDir, ".pwgen.yaml"),
			filepath.Join(homeDir, ".pwgen.yml"),
			filepath.Join(homeDir, ".config", "pwgen", "config.yaml"),
			filepath.Join(homeDir, ".config", "pwgen", "config.yml"),
		)
	}

	for _, path := range configPaths {
		if err := loadConfigFromFile(path, &config); err == nil {
			break // Use first config file found
		}
	}

	// Override with environment variables
	loadConfigFromEnv(&config)

	return config, nil
}

func loadConfigFromFile(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func loadConfigFromEnv(config *Config) {
	if val := os.Getenv("PWGEN_LENGTH"); val != "" {
		if length, err := strconv.Atoi(val); err == nil {
			config.Length = length
		}
	}

	if val := os.Getenv("PWGEN_INCLUDE_UPPER"); val != "" {
		config.IncludeUpper = parseBool(val, config.IncludeUpper)
	}

	if val := os.Getenv("PWGEN_INCLUDE_LOWER"); val != "" {
		config.IncludeLower = parseBool(val, config.IncludeLower)
	}

	if val := os.Getenv("PWGEN_INCLUDE_DIGITS"); val != "" {
		config.IncludeDigits = parseBool(val, config.IncludeDigits)
	}

	if val := os.Getenv("PWGEN_INCLUDE_SYMBOLS"); val != "" {
		config.IncludeSymbols = parseBool(val, config.IncludeSymbols)
	}

	if val := os.Getenv("PWGEN_EXCLUDE_AMBIGUOUS"); val != "" {
		config.ExcludeAmbiguous = parseBool(val, config.ExcludeAmbiguous)
	}

	if val := os.Getenv("PWGEN_COUNT"); val != "" {
		if count, err := strconv.Atoi(val); err == nil {
			config.Count = count
		}
	}

	if val := os.Getenv("PWGEN_SHOW_STRENGTH"); val != "" {
		config.ShowStrength = parseBool(val, config.ShowStrength)
	}

	if val := os.Getenv("PWGEN_POLICY_TEMPLATE"); val != "" {
		config.PolicyTemplate = val
	}
}

func parseBool(val string, defaultValue bool) bool {
	switch strings.ToLower(val) {
	case "true", "1", "yes", "on", "enable", "enabled":
		return true
	case "false", "0", "no", "off", "disable", "disabled":
		return false
	default:
		return defaultValue
	}
}

func (c Config) ToPasswordConfig() PasswordConfig {
	return PasswordConfig{
		Length:           c.Length,
		IncludeUpper:     c.IncludeUpper,
		IncludeLower:     c.IncludeLower,
		IncludeDigits:    c.IncludeDigits,
		IncludeSymbols:   c.IncludeSymbols,
		ExcludeAmbiguous: c.ExcludeAmbiguous,
	}
}

func SaveConfigExample(path string) error {
	config := Config{
		Length:           16,
		IncludeUpper:     true,
		IncludeLower:     true,
		IncludeDigits:    true,
		IncludeSymbols:   true,
		ExcludeAmbiguous: true,
		Count:            1,
		ShowStrength:     true,
		PolicyTemplate:   "corporate",
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := `# Password Generator Configuration
# This file contains default settings for the password generator
# Place this file in your home directory as ~/.pwgen.yaml or in the current directory as .pwgen.yaml
# Environment variables (PWGEN_*) will override these settings
# Command-line flags will override both config file and environment variables

`

	content := header + string(data)

	return os.WriteFile(path, []byte(content), 0644)
}
