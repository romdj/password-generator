package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

type PasswordConfig struct {
	Length      int
	IncludeUpper bool
	IncludeLower bool
	IncludeDigits bool
	IncludeSymbols bool
	ExcludeAmbiguous bool
}

const (
	LowerCase = "abcdefghijklmnopqrstuvwxyz"
	UpperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits = "0123456789"
	Symbols = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	Ambiguous = "0O1lI"
)

func main() {
	// Load configuration from files and environment
	baseConfig, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load config: %v\n", err)
		baseConfig = DefaultConfig()
	}
	
	// Convert to PasswordConfig for compatibility
	config := baseConfig.ToPasswordConfig()
	count := baseConfig.Count
	showStrength := baseConfig.ShowStrength
	policyTemplate := baseConfig.PolicyTemplate
	
	// Command line flags override config
	flag.IntVar(&config.Length, "length", config.Length, "Password length")
	flag.IntVar(&config.Length, "l", config.Length, "Password length (short)")
	flag.BoolVar(&config.IncludeUpper, "upper", config.IncludeUpper, "Include uppercase letters")
	flag.BoolVar(&config.IncludeUpper, "u", config.IncludeUpper, "Include uppercase letters (short)")
	flag.BoolVar(&config.IncludeLower, "lower", config.IncludeLower, "Include lowercase letters")
	flag.BoolVar(&config.IncludeLower, "L", config.IncludeLower, "Include lowercase letters (short)")
	flag.BoolVar(&config.IncludeDigits, "digits", config.IncludeDigits, "Include digits")
	flag.BoolVar(&config.IncludeDigits, "d", config.IncludeDigits, "Include digits (short)")
	flag.BoolVar(&config.IncludeSymbols, "symbols", config.IncludeSymbols, "Include symbols")
	flag.BoolVar(&config.IncludeSymbols, "s", config.IncludeSymbols, "Include symbols (short)")
	flag.BoolVar(&config.ExcludeAmbiguous, "no-ambiguous", config.ExcludeAmbiguous, "Exclude ambiguous characters (0, O, 1, l, I)")
	flag.BoolVar(&config.ExcludeAmbiguous, "n", config.ExcludeAmbiguous, "Exclude ambiguous characters (short)")
	
	flag.IntVar(&count, "count", count, "Number of passwords to generate")
	countShort := flag.Int("c", count, "Number of passwords to generate (short)")
	flag.BoolVar(&showStrength, "strength", showStrength, "Show password strength analysis")
	flag.BoolVar(&showStrength, "S", showStrength, "Show password strength analysis (short)")
	flag.StringVar(&policyTemplate, "policy", policyTemplate, "Apply password policy template")
	flag.StringVar(&policyTemplate, "p", policyTemplate, "Apply password policy template (short)")
	
	listPolicies := flag.Bool("list-policies", false, "List available password policy templates")
	validateOnly := flag.String("validate", "", "Validate a password against policy without generating")
	saveConfig := flag.String("save-config", "", "Save example configuration to file")
	
	flag.Parse()
	
	// Handle special commands
	if *listPolicies {
		fmt.Println("Available password policy templates:")
		for _, name := range ListPolicies() {
			policy, _ := GetPolicy(name)
			fmt.Printf("  %-15s - %s\n", name, policy.Description)
		}
		return
	}
	
	if *saveConfig != "" {
		if err := SaveConfigExample(*saveConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Example configuration saved to %s\n", *saveConfig)
		return
	}
	
	if *validateOnly != "" {
		if policyTemplate == "" {
			fmt.Fprintf(os.Stderr, "Error: --policy required when using --validate\n")
			os.Exit(1)
		}
		
		policy, err := GetPolicy(policyTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		
		violations := ValidatePasswordAgainstPolicy(*validateOnly, policy)
		if len(violations) == 0 {
			fmt.Printf("✓ Password meets %s policy requirements\n", policy.Name)
		} else {
			fmt.Printf("✗ Password violates %s policy:\n", policy.Name)
			for _, violation := range violations {
				fmt.Printf("  - %s\n", violation.Description)
			}
		}
		return
	}
	
	// Apply policy template if specified
	var policy PasswordPolicy
	if policyTemplate != "" {
		p, err := GetPolicy(policyTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintf(os.Stderr, "Available policies: %s\n", strings.Join(ListPolicies(), ", "))
			os.Exit(1)
		}
		policy = p
		ApplyPolicyToConfig(policy, &config)
	}
	
	// Use short flag if set
	if *countShort != count {
		count = *countShort
	}
	
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	for i := 0; i < count; i++ {
		password, err := generatePassword(config)
		if err != nil {
			log.Fatalf("Failed to generate password: %v", err)
		}
		
		fmt.Print(password)
		
		// Show strength analysis if requested
		if showStrength {
			strength := AnalyzePasswordStrength(password)
			fmt.Printf(" [%s%s\033[0m, Score: %d/100, Entropy: %.1f bits, Time to crack: %s]",
				strength.Level.Color(),
				strength.Level.String(),
				strength.Score,
				strength.Entropy,
				strength.TimeToCrack,
			)
			
			if len(strength.Feedback) > 0 {
				fmt.Printf("\n  Feedback: %s", strings.Join(strength.Feedback, "; "))
			}
		}
		
		// Validate against policy if specified
		if policyTemplate != "" {
			violations := ValidatePasswordAgainstPolicy(password, policy)
			if len(violations) > 0 {
				fmt.Printf(" [Policy violations: %d]", len(violations))
				if showStrength {
					fmt.Printf("\n  Violations:")
					for _, violation := range violations {
						fmt.Printf("\n    - %s", violation.Description)
					}
				}
			}
		}
		
		fmt.Println()
	}
}

func validateConfig(config PasswordConfig) error {
	if config.Length < 1 {
		return fmt.Errorf("password length must be at least 1")
	}
	
	if !config.IncludeUpper && !config.IncludeLower && !config.IncludeDigits && !config.IncludeSymbols {
		return fmt.Errorf("at least one character type must be enabled")
	}
	
	return nil
}

func generatePassword(config PasswordConfig) (string, error) {
	charset := buildCharset(config)
	
	if len(charset) == 0 {
		return "", fmt.Errorf("no valid characters available for password generation")
	}
	
	password := make([]byte, config.Length)
	
	for i := 0; i < config.Length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = charset[randomIndex.Int64()]
	}
	
	return string(password), nil
}

func buildCharset(config PasswordConfig) string {
	var charset strings.Builder
	
	if config.IncludeLower {
		charset.WriteString(LowerCase)
	}
	
	if config.IncludeUpper {
		charset.WriteString(UpperCase)
	}
	
	if config.IncludeDigits {
		charset.WriteString(Digits)
	}
	
	if config.IncludeSymbols {
		charset.WriteString(Symbols)
	}
	
	result := charset.String()
	
	if config.ExcludeAmbiguous {
		for _, char := range Ambiguous {
			result = strings.ReplaceAll(result, string(char), "")
		}
	}
	
	return result
}