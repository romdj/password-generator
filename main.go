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
	config := PasswordConfig{}
	
	flag.IntVar(&config.Length, "length", 12, "Password length")
	flag.IntVar(&config.Length, "l", 12, "Password length (short)")
	flag.BoolVar(&config.IncludeUpper, "upper", true, "Include uppercase letters")
	flag.BoolVar(&config.IncludeUpper, "u", true, "Include uppercase letters (short)")
	flag.BoolVar(&config.IncludeLower, "lower", true, "Include lowercase letters")
	flag.BoolVar(&config.IncludeLower, "L", true, "Include lowercase letters (short)")
	flag.BoolVar(&config.IncludeDigits, "digits", true, "Include digits")
	flag.BoolVar(&config.IncludeDigits, "d", true, "Include digits (short)")
	flag.BoolVar(&config.IncludeSymbols, "symbols", false, "Include symbols")
	flag.BoolVar(&config.IncludeSymbols, "s", false, "Include symbols (short)")
	flag.BoolVar(&config.ExcludeAmbiguous, "no-ambiguous", false, "Exclude ambiguous characters (0, O, 1, l, I)")
	flag.BoolVar(&config.ExcludeAmbiguous, "n", false, "Exclude ambiguous characters (short)")
	
	count := flag.Int("count", 1, "Number of passwords to generate")
	count_short := flag.Int("c", 1, "Number of passwords to generate (short)")
	
	flag.Parse()
	
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	passwordCount := *count
	if *count_short != 1 {
		passwordCount = *count_short
	}
	
	for i := 0; i < passwordCount; i++ {
		password, err := generatePassword(config)
		if err != nil {
			log.Fatalf("Failed to generate password: %v", err)
		}
		fmt.Println(password)
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