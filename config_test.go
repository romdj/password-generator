package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Length != 12 {
		t.Errorf("DefaultConfig() Length = %d, want 12", config.Length)
	}

	if !config.IncludeUpper {
		t.Error("DefaultConfig() should include uppercase by default")
	}

	if !config.IncludeLower {
		t.Error("DefaultConfig() should include lowercase by default")
	}

	if !config.IncludeDigits {
		t.Error("DefaultConfig() should include digits by default")
	}

	if config.IncludeSymbols {
		t.Error("DefaultConfig() should not include symbols by default")
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set test environment variables
	os.Setenv("PWGEN_LENGTH", "16")
	os.Setenv("PWGEN_INCLUDE_SYMBOLS", "true")
	os.Setenv("PWGEN_SHOW_STRENGTH", "yes")
	defer func() {
		os.Unsetenv("PWGEN_LENGTH")
		os.Unsetenv("PWGEN_INCLUDE_SYMBOLS")
		os.Unsetenv("PWGEN_SHOW_STRENGTH")
	}()

	config := DefaultConfig()
	loadConfigFromEnv(&config)

	if config.Length != 16 {
		t.Errorf("loadConfigFromEnv() Length = %d, want 16", config.Length)
	}

	if !config.IncludeSymbols {
		t.Error("loadConfigFromEnv() should enable symbols")
	}

	if !config.ShowStrength {
		t.Error("loadConfigFromEnv() should enable strength display")
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue bool
		want         bool
	}{
		{"true", "true", false, true},
		{"1", "1", false, true},
		{"yes", "yes", false, true},
		{"enabled", "enabled", false, true},
		{"false", "false", true, false},
		{"0", "0", true, false},
		{"no", "no", true, false},
		{"disabled", "disabled", true, false},
		{"invalid", "invalid", true, true},
		{"empty", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBool(tt.value, tt.defaultValue)
			if got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigToPasswordConfig(t *testing.T) {
	config := Config{
		Length:           16,
		IncludeUpper:     true,
		IncludeLower:     false,
		IncludeDigits:    true,
		IncludeSymbols:   true,
		ExcludeAmbiguous: true,
	}

	pwConfig := config.ToPasswordConfig()

	if pwConfig.Length != config.Length {
		t.Errorf("ToPasswordConfig() Length = %d, want %d", pwConfig.Length, config.Length)
	}

	if pwConfig.IncludeUpper != config.IncludeUpper {
		t.Errorf("ToPasswordConfig() IncludeUpper = %v, want %v", pwConfig.IncludeUpper, config.IncludeUpper)
	}

	if pwConfig.IncludeLower != config.IncludeLower {
		t.Errorf("ToPasswordConfig() IncludeLower = %v, want %v", pwConfig.IncludeLower, config.IncludeLower)
	}
}

func TestLoadConfig(t *testing.T) {
	// Test loading config when no config files exist
	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
	}

	// Should return default config when no files exist
	defaultConfig := DefaultConfig()
	if config.Length != defaultConfig.Length {
		t.Errorf("LoadConfig() Length = %d, want %d", config.Length, defaultConfig.Length)
	}
}

func TestLoadConfigWithFile(t *testing.T) {
	// Save current directory to restore later
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary directory and change to it
	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	defer func() {
		os.Chdir(originalDir)
	}()

	// Create a test config file
	configContent := `length: 20
include_upper: false
include_symbols: true
count: 3`

	err = os.WriteFile(".pwgen.yaml", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
	}

	if config.Length != 20 {
		t.Errorf("LoadConfig() Length = %d, want 20", config.Length)
	}

	if config.IncludeUpper != false {
		t.Errorf("LoadConfig() IncludeUpper = %v, want false", config.IncludeUpper)
	}

	if config.IncludeSymbols != true {
		t.Errorf("LoadConfig() IncludeSymbols = %v, want true", config.IncludeSymbols)
	}

	if config.Count != 3 {
		t.Errorf("LoadConfig() Count = %d, want 3", config.Count)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	// Create a test config file
	configContent := `length: 20
include_upper: false
include_symbols: true
count: 3`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config := DefaultConfig()
	err = loadConfigFromFile(configPath, &config)
	if err != nil {
		t.Errorf("loadConfigFromFile() error = %v", err)
	}

	if config.Length != 20 {
		t.Errorf("loadConfigFromFile() Length = %d, want 20", config.Length)
	}

	if config.IncludeUpper != false {
		t.Errorf("loadConfigFromFile() IncludeUpper = %v, want false", config.IncludeUpper)
	}

	if config.IncludeSymbols != true {
		t.Errorf("loadConfigFromFile() IncludeSymbols = %v, want true", config.IncludeSymbols)
	}

	// Test with non-existent file
	err = loadConfigFromFile("nonexistent.yaml", &config)
	if err == nil {
		t.Error("loadConfigFromFile() should return error for non-existent file")
	}
}

func TestLoadConfigFromEnvExtended(t *testing.T) {
	// Test all environment variables
	envVars := map[string]string{
		"PWGEN_LENGTH":            "24",
		"PWGEN_INCLUDE_UPPER":     "false",
		"PWGEN_INCLUDE_LOWER":     "true",
		"PWGEN_INCLUDE_DIGITS":    "false",
		"PWGEN_INCLUDE_SYMBOLS":   "true",
		"PWGEN_EXCLUDE_AMBIGUOUS": "true",
		"PWGEN_COUNT":             "5",
		"PWGEN_SHOW_STRENGTH":     "true",
		"PWGEN_POLICY_TEMPLATE":   "high-security",
	}

	// Set environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	config := DefaultConfig()
	loadConfigFromEnv(&config)

	if config.Length != 24 {
		t.Errorf("loadConfigFromEnv() Length = %d, want 24", config.Length)
	}

	if config.IncludeUpper != false {
		t.Errorf("loadConfigFromEnv() IncludeUpper = %v, want false", config.IncludeUpper)
	}

	if config.Count != 5 {
		t.Errorf("loadConfigFromEnv() Count = %d, want 5", config.Count)
	}

	if config.PolicyTemplate != "high-security" {
		t.Errorf("loadConfigFromEnv() PolicyTemplate = %s, want high-security", config.PolicyTemplate)
	}
}

func TestSaveConfigExample(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	err := SaveConfigExample(configPath)
	if err != nil {
		t.Errorf("SaveConfigExample() error = %v", err)
		return
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("SaveConfigExample() did not create config file")
		return
	}

	// Read and verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Errorf("Failed to read created config file: %v", err)
		return
	}

	contentStr := string(content)
	if len(contentStr) == 0 {
		t.Error("SaveConfigExample() created empty config file")
	}

	// Check for expected content
	expectedStrings := []string{
		"length:",
		"include_upper:",
		"include_symbols:",
		"policy_template:",
	}

	for _, expected := range expectedStrings {
		if !contains(contentStr, expected) {
			t.Errorf("SaveConfigExample() missing expected content: %s", expected)
		}
	}
}

func TestSaveConfigExampleError(t *testing.T) {
	// Test error case - try to write to invalid path
	invalidPath := "/dev/null/invalid/path/config.yaml"
	err := SaveConfigExample(invalidPath)
	if err == nil {
		t.Error("SaveConfigExample() should return error for invalid path")
	}
}

func TestSaveConfigExampleValidPath(t *testing.T) {
	// Test normal operation to improve coverage
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "valid-config.yaml")

	err := SaveConfigExample(configPath)
	if err != nil {
		t.Errorf("SaveConfigExample() error = %v", err)
	}

	// Verify file was created and has content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Errorf("Failed to read created config file: %v", err)
	}

	if len(content) == 0 {
		t.Error("SaveConfigExample() created empty config file")
	}
}

// Note: Testing yaml.Marshal error path is difficult without complex mocking
// as yaml.Marshal rarely fails with normal struct data. In practice, this
// error path would only be hit with corrupted memory or similar system issues.

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			contains(s[1:], substr))))
}
