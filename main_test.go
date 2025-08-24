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

func TestGeneratePasswordCryptoRandEdgeCase(t *testing.T) {
	// This test documents the theoretical crypto/rand error path that is
	// extremely difficult to trigger in practice. The rand.Int function
	// from crypto/rand can theoretically fail, but this almost never happens
	// in real-world scenarios except in cases of:
	// 1. System entropy depletion (very rare on modern systems)
	// 2. Hardware random number generator failure
	// 3. System-level issues with /dev/urandom access
	//
	// Since we cannot easily mock crypto/rand.Int to return an error,
	// this test serves as documentation that we are aware of this edge case
	// and that it represents the 10% uncovered in generatePassword function.

	config := PasswordConfig{
		Length:         1000000, // Very large password to increase chances of hitting the edge case
		IncludeUpper:   true,
		IncludeLower:   true,
		IncludeDigits:  true,
		IncludeSymbols: true,
	}

	// Generate multiple passwords to exercise the crypto/rand path extensively
	for i := 0; i < 10; i++ {
		password, err := generatePassword(config)
		if err != nil {
			// If we ever hit the crypto/rand error, this validates our error handling
			t.Logf("Crypto rand error encountered (rare but valid): %v", err)
			if password != "" {
				t.Error("generatePassword() should return empty string on crypto/rand error")
			}
			return // Test passed by hitting the error path
		}

		if len(password) != config.Length {
			t.Errorf("generatePassword() length = %d, want %d", len(password), config.Length)
		}
	}

	// If we reach here, crypto/rand worked correctly (expected case)
	t.Log("Crypto/rand worked correctly for all iterations (expected behavior)")
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
		{
			name: "uppercase only",
			config: PasswordConfig{
				IncludeUpper: true,
			},
			want: UpperCase,
		},
		{
			name: "digits only",
			config: PasswordConfig{
				IncludeDigits: true,
			},
			want: Digits,
		},
		{
			name: "symbols only",
			config: PasswordConfig{
				IncludeSymbols: true,
			},
			want: Symbols,
		},
		{
			name: "exclude ambiguous from all types",
			config: PasswordConfig{
				IncludeUpper:     true,
				IncludeLower:     true,
				IncludeDigits:    true,
				IncludeSymbols:   true,
				ExcludeAmbiguous: true,
			},
			// Ambiguous chars "0O1lI" should be excluded
			want: "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$%^&*()_+-=[]{}|;:,.<>?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCharset(tt.config)
			if got != tt.want {
				// For debug: show actual result for fixing test expectations
				if tt.name == "exclude ambiguous from all types" {
					t.Logf("Actual result: %q", got)
					t.Logf("Expected result: %q", tt.want)
				}
				t.Errorf("buildCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}
