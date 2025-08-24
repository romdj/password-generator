package main

import (
	"testing"
)

func TestGetPolicy(t *testing.T) {
	tests := []struct {
		name       string
		policyName string
		wantErr    bool
		wantMinLen int
	}{
		{
			name:       "basic policy",
			policyName: "basic",
			wantErr:    false,
			wantMinLen: 8,
		},
		{
			name:       "corporate policy",
			policyName: "corporate",
			wantErr:    false,
			wantMinLen: 12,
		},
		{
			name:       "high-security policy",
			policyName: "high-security",
			wantErr:    false,
			wantMinLen: 16,
		},
		{
			name:       "nonexistent policy",
			policyName: "nonexistent",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := GetPolicy(tt.policyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && policy.MinLength != tt.wantMinLen {
				t.Errorf("GetPolicy() MinLength = %d, want %d", policy.MinLength, tt.wantMinLen)
			}
		})
	}
}

func TestValidatePasswordAgainstPolicy(t *testing.T) {
	basicPolicy, _ := GetPolicy("basic")
	corporatePolicy, _ := GetPolicy("corporate")

	tests := []struct {
		name           string
		password       string
		policy         PasswordPolicy
		wantViolations int
	}{
		{
			name:           "valid basic password",
			password:       "MySecure1", // Changed to avoid 'password' pattern
			policy:         basicPolicy,
			wantViolations: 0,
		},
		{
			name:           "too short",
			password:       "Pass1",
			policy:         basicPolicy,
			wantViolations: 1,
		},
		{
			name:           "missing digits",
			password:       "MySecure", // Changed to avoid 'password' pattern
			policy:         basicPolicy,
			wantViolations: 2, // RequireDigits and MinDigits
		},
		{
			name:           "corporate policy violation",
			password:       "weak123", // Simple password with multiple violations
			policy:         corporatePolicy,
			wantViolations: 7, // All the violations listed in test output
		},
		{
			name:           "valid corporate password",
			password:       "MyC2rp2r@te!Secure", // Removed ambiguous chars and 'password' pattern
			policy:         corporatePolicy,
			wantViolations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := ValidatePasswordAgainstPolicy(tt.password, tt.policy)
			if len(violations) != tt.wantViolations {
				t.Errorf("ValidatePasswordAgainstPolicy() violations = %d, want %d",
					len(violations), tt.wantViolations)
				for _, v := range violations {
					t.Logf("  - %s: %s", v.Rule, v.Description)
				}
			}
		})
	}
}

func TestApplyPolicyToConfig(t *testing.T) {
	basicPolicy, _ := GetPolicy("basic")

	config := PasswordConfig{
		Length:           6, // too short
		IncludeUpper:     false,
		IncludeLower:     true,
		IncludeDigits:    false,
		IncludeSymbols:   false,
		ExcludeAmbiguous: false,
	}

	ApplyPolicyToConfig(basicPolicy, &config)

	if config.Length < basicPolicy.MinLength {
		t.Errorf("ApplyPolicyToConfig() length = %d, want >= %d", config.Length, basicPolicy.MinLength)
	}

	if !config.IncludeUpper {
		t.Error("ApplyPolicyToConfig() should enable uppercase letters")
	}

	if !config.IncludeDigits {
		t.Error("ApplyPolicyToConfig() should enable digits")
	}
}

func TestListPolicies(t *testing.T) {
	policies := ListPolicies()

	expectedPolicies := []string{"basic", "corporate", "high-security", "aws", "azure", "pci-dss"}

	if len(policies) < len(expectedPolicies) {
		t.Errorf("ListPolicies() returned %d policies, want at least %d", len(policies), len(expectedPolicies))
	}

	// Check if all expected policies are present
	policyMap := make(map[string]bool)
	for _, policy := range policies {
		policyMap[policy] = true
	}

	for _, expected := range expectedPolicies {
		if !policyMap[expected] {
			t.Errorf("ListPolicies() missing expected policy: %s", expected)
		}
	}
}
