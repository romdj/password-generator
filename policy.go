package main

import (
	"fmt"
	"regexp"
	"strings"
)

type PasswordPolicy struct {
	Name              string   `yaml:"name"`
	Description       string   `yaml:"description"`
	MinLength         int      `yaml:"min_length"`
	MaxLength         int      `yaml:"max_length"`
	RequireUpper      bool     `yaml:"require_upper"`
	RequireLower      bool     `yaml:"require_lower"`
	RequireDigits     bool     `yaml:"require_digits"`
	RequireSymbols    bool     `yaml:"require_symbols"`
	MinUpper          int      `yaml:"min_upper"`
	MinLower          int      `yaml:"min_lower"`
	MinDigits         int      `yaml:"min_digits"`
	MinSymbols        int      `yaml:"min_symbols"`
	ExcludeAmbiguous  bool     `yaml:"exclude_ambiguous"`
	ForbiddenChars    string   `yaml:"forbidden_chars"`
	ForbiddenPatterns []string `yaml:"forbidden_patterns"`
	MinEntropy        float64  `yaml:"min_entropy"`
}

type PolicyViolation struct {
	Rule        string
	Description string
}

var BuiltinPolicies = map[string]PasswordPolicy{
	"basic": {
		Name:              "Basic Security",
		Description:       "Minimal security requirements suitable for low-risk applications",
		MinLength:         8,
		MaxLength:         0, // No limit
		RequireUpper:      true,
		RequireLower:      true,
		RequireDigits:     true,
		RequireSymbols:    false,
		MinUpper:          1,
		MinLower:          1,
		MinDigits:         1,
		MinSymbols:        0,
		ExcludeAmbiguous:  false,
		ForbiddenChars:    "",
		ForbiddenPatterns: []string{"password", "123456", "qwerty"},
		MinEntropy:        25,
	},
	"corporate": {
		Name:              "Corporate Standard",
		Description:       "Standard corporate password policy with moderate security",
		MinLength:         12,
		MaxLength:         0,
		RequireUpper:      true,
		RequireLower:      true,
		RequireDigits:     true,
		RequireSymbols:    true,
		MinUpper:          2,
		MinLower:          2,
		MinDigits:         2,
		MinSymbols:        1,
		ExcludeAmbiguous:  true,
		ForbiddenChars:    "",
		ForbiddenPatterns: []string{"password", "123456", "qwerty", "admin", "login", "welcome"},
		MinEntropy:        40,
	},
	"high-security": {
		Name:             "High Security",
		Description:      "Stringent requirements for high-security environments",
		MinLength:        16,
		MaxLength:        0,
		RequireUpper:     true,
		RequireLower:     true,
		RequireDigits:    true,
		RequireSymbols:   true,
		MinUpper:         3,
		MinLower:         3,
		MinDigits:        3,
		MinSymbols:       2,
		ExcludeAmbiguous: true,
		ForbiddenChars:   "",
		ForbiddenPatterns: []string{
			"password", "123456", "qwerty", "admin", "login", "welcome",
			"letmein", "monkey", "dragon", "master", "shadow", "football",
		},
		MinEntropy: 60,
	},
	"aws": {
		Name:              "AWS IAM Policy",
		Description:       "Meets AWS IAM password policy requirements",
		MinLength:         8,
		MaxLength:         128,
		RequireUpper:      true,
		RequireLower:      true,
		RequireDigits:     true,
		RequireSymbols:    true,
		MinUpper:          1,
		MinLower:          1,
		MinDigits:         1,
		MinSymbols:        1,
		ExcludeAmbiguous:  false,
		ForbiddenChars:    "",
		ForbiddenPatterns: []string{},
		MinEntropy:        30,
	},
	"azure": {
		Name:              "Azure AD Policy",
		Description:       "Meets Azure AD password complexity requirements",
		MinLength:         8,
		MaxLength:         256,
		RequireUpper:      true,
		RequireLower:      true,
		RequireDigits:     true,
		RequireSymbols:    true,
		MinUpper:          1,
		MinLower:          1,
		MinDigits:         1,
		MinSymbols:        1,
		ExcludeAmbiguous:  false,
		ForbiddenChars:    "",
		ForbiddenPatterns: []string{},
		MinEntropy:        35,
	},
	"pci-dss": {
		Name:              "PCI DSS Compliant",
		Description:       "Meets PCI DSS password requirements for payment systems",
		MinLength:         7,
		MaxLength:         0,
		RequireUpper:      true,
		RequireLower:      true,
		RequireDigits:     true,
		RequireSymbols:    false,
		MinUpper:          1,
		MinLower:          1,
		MinDigits:         1,
		MinSymbols:        0,
		ExcludeAmbiguous:  false,
		ForbiddenChars:    "",
		ForbiddenPatterns: []string{},
		MinEntropy:        28,
	},
}

func GetPolicy(name string) (PasswordPolicy, error) {
	if policy, exists := BuiltinPolicies[name]; exists {
		return policy, nil
	}
	return PasswordPolicy{}, fmt.Errorf("policy '%s' not found", name)
}

func ListPolicies() []string {
	var policies []string
	for name := range BuiltinPolicies {
		policies = append(policies, name)
	}
	return policies
}

func ValidatePasswordAgainstPolicy(password string, policy PasswordPolicy) []PolicyViolation {
	var violations []PolicyViolation

	// Length checks
	if len(password) < policy.MinLength {
		violations = append(violations, PolicyViolation{
			Rule:        "MinLength",
			Description: fmt.Sprintf("Password must be at least %d characters long", policy.MinLength),
		})
	}

	if policy.MaxLength > 0 && len(password) > policy.MaxLength {
		violations = append(violations, PolicyViolation{
			Rule:        "MaxLength",
			Description: fmt.Sprintf("Password must not exceed %d characters", policy.MaxLength),
		})
	}

	// Character type requirements
	upperCount := countMatches(password, `[A-Z]`)
	lowerCount := countMatches(password, `[a-z]`)
	digitCount := countMatches(password, `[0-9]`)
	symbolCount := countMatches(password, `[^a-zA-Z0-9]`)

	if policy.RequireUpper && upperCount == 0 {
		violations = append(violations, PolicyViolation{
			Rule:        "RequireUpper",
			Description: "Password must contain at least one uppercase letter",
		})
	}

	if policy.RequireLower && lowerCount == 0 {
		violations = append(violations, PolicyViolation{
			Rule:        "RequireLower",
			Description: "Password must contain at least one lowercase letter",
		})
	}

	if policy.RequireDigits && digitCount == 0 {
		violations = append(violations, PolicyViolation{
			Rule:        "RequireDigits",
			Description: "Password must contain at least one digit",
		})
	}

	if policy.RequireSymbols && symbolCount == 0 {
		violations = append(violations, PolicyViolation{
			Rule:        "RequireSymbols",
			Description: "Password must contain at least one symbol",
		})
	}

	// Minimum character counts
	if upperCount < policy.MinUpper {
		violations = append(violations, PolicyViolation{
			Rule:        "MinUpper",
			Description: fmt.Sprintf("Password must contain at least %d uppercase letters", policy.MinUpper),
		})
	}

	if lowerCount < policy.MinLower {
		violations = append(violations, PolicyViolation{
			Rule:        "MinLower",
			Description: fmt.Sprintf("Password must contain at least %d lowercase letters", policy.MinLower),
		})
	}

	if digitCount < policy.MinDigits {
		violations = append(violations, PolicyViolation{
			Rule:        "MinDigits",
			Description: fmt.Sprintf("Password must contain at least %d digits", policy.MinDigits),
		})
	}

	if symbolCount < policy.MinSymbols {
		violations = append(violations, PolicyViolation{
			Rule:        "MinSymbols",
			Description: fmt.Sprintf("Password must contain at least %d symbols", policy.MinSymbols),
		})
	}

	// Ambiguous character check
	if policy.ExcludeAmbiguous {
		ambiguous := "0O1lI"
		for _, char := range ambiguous {
			if strings.ContainsRune(password, char) {
				violations = append(violations, PolicyViolation{
					Rule:        "ExcludeAmbiguous",
					Description: fmt.Sprintf("Password must not contain ambiguous characters (%s)", ambiguous),
				})
				break
			}
		}
	}

	// Forbidden characters
	if policy.ForbiddenChars != "" {
		for _, char := range policy.ForbiddenChars {
			if strings.ContainsRune(password, char) {
				violations = append(violations, PolicyViolation{
					Rule:        "ForbiddenChars",
					Description: fmt.Sprintf("Password must not contain forbidden character '%c'", char),
				})
			}
		}
	}

	// Forbidden patterns
	lower := strings.ToLower(password)
	for _, pattern := range policy.ForbiddenPatterns {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			violations = append(violations, PolicyViolation{
				Rule:        "ForbiddenPatterns",
				Description: fmt.Sprintf("Password must not contain forbidden pattern '%s'", pattern),
			})
		}
	}

	// Entropy check
	if policy.MinEntropy > 0 {
		entropy := calculateEntropy(password)
		if entropy < policy.MinEntropy {
			violations = append(violations, PolicyViolation{
				Rule:        "MinEntropy",
				Description: fmt.Sprintf("Password entropy (%.1f bits) must be at least %.1f bits", entropy, policy.MinEntropy),
			})
		}
	}

	return violations
}

func countMatches(text, pattern string) int {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)
	return len(matches)
}

func ApplyPolicyToConfig(policy PasswordPolicy, config *PasswordConfig) {
	// Adjust length to meet minimum requirements
	if config.Length < policy.MinLength {
		config.Length = policy.MinLength
	}

	if policy.MaxLength > 0 && config.Length > policy.MaxLength {
		config.Length = policy.MaxLength
	}

	// Enable required character types
	if policy.RequireUpper {
		config.IncludeUpper = true
	}

	if policy.RequireLower {
		config.IncludeLower = true
	}

	if policy.RequireDigits {
		config.IncludeDigits = true
	}

	if policy.RequireSymbols {
		config.IncludeSymbols = true
	}

	// Apply ambiguous character exclusion
	if policy.ExcludeAmbiguous {
		config.ExcludeAmbiguous = true
	}
}
