package main

import (
	"strings"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name   string
		config PasswordConfig
		want   int // expected length
	}{
		{
			name: "basic password",
			config: PasswordConfig{
				Length:        12,
				IncludeUpper:  true,
				IncludeLower:  true,
				IncludeDigits: true,
			},
			want: 12,
		},
		{
			name: "symbols only",
			config: PasswordConfig{
				Length:         8,
				IncludeSymbols: true,
			},
			want: 8,
		},
		{
			name: "exclude ambiguous",
			config: PasswordConfig{
				Length:           16,
				IncludeUpper:     true,
				IncludeLower:     true,
				IncludeDigits:    true,
				ExcludeAmbiguous: true,
			},
			want: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := generatePassword(tt.config)
			if err != nil {
				t.Errorf("generatePassword() error = %v", err)
				return
			}

			if len(password) != tt.want {
				t.Errorf("generatePassword() length = %v, want %v", len(password), tt.want)
			}

			// Test ambiguous character exclusion
			if tt.config.ExcludeAmbiguous {
				for _, char := range Ambiguous {
					if strings.ContainsRune(password, char) {
						t.Errorf("generatePassword() contains ambiguous character %c", char)
					}
				}
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  PasswordConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: PasswordConfig{
				Length:       12,
				IncludeUpper: true,
			},
			wantErr: false,
		},
		{
			name: "zero length",
			config: PasswordConfig{
				Length: 0,
			},
			wantErr: true,
		},
		{
			name: "no character types",
			config: PasswordConfig{
				Length: 12,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeneratePasswordErrorCases(t *testing.T) {
	// Test error case where no character types are enabled
	config := PasswordConfig{
		Length: 10,
		// All character types disabled
	}

	password, err := generatePassword(config)
	if err == nil {
		t.Error("generatePassword() should return error when no character types enabled")
	}

	if password != "" {
		t.Errorf("generatePassword() should return empty string on error, got %s", password)
	}
}

func TestGeneratePasswordEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		config  PasswordConfig
		wantErr bool
		wantLen int
	}{
		{
			name: "zero length password",
			config: PasswordConfig{
				Length:       0,
				IncludeUpper: true,
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "single character password",
			config: PasswordConfig{
				Length:       1,
				IncludeUpper: true,
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "very long password",
			config: PasswordConfig{
				Length:         100,
				IncludeUpper:   true,
				IncludeLower:   true,
				IncludeDigits:  true,
				IncludeSymbols: true,
			},
			wantErr: false,
			wantLen: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := generatePassword(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("generatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(password) != tt.wantLen {
				t.Errorf("generatePassword() length = %d, want %d", len(password), tt.wantLen)
			}
		})
	}
}

func TestBuildCharset(t *testing.T) {
	tests := []struct {
		name   string
		config PasswordConfig
		want   string
	}{
		{
			name: "all character types",
			config: PasswordConfig{
				IncludeUpper:   true,
				IncludeLower:   true,
				IncludeDigits:  true,
				IncludeSymbols: true,
			},
			want: LowerCase + UpperCase + Digits + Symbols,
		},
		{
			name: "lowercase only",
			config: PasswordConfig{
				IncludeLower: true,
			},
			want: LowerCase,
		},
		{
			name: "exclude ambiguous",
			config: PasswordConfig{
				IncludeLower:     true,
				IncludeDigits:    true,
				ExcludeAmbiguous: true,
			},
			want: "abcdefghijkmnopqrstuvwxyz23456789", // Excludes l, 1, 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCharset(tt.config)
			if got != tt.want {
				t.Errorf("buildCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}
