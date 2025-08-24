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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		contains(s[1:], substr))))
}