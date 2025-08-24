package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type StrengthLevel int

const (
	VeryWeak StrengthLevel = iota
	Weak
	Fair
	Good
	Strong
	VeryStrong
)

func (s StrengthLevel) String() string {
	switch s {
	case VeryWeak:
		return "Very Weak"
	case Weak:
		return "Weak"
	case Fair:
		return "Fair"
	case Good:
		return "Good"
	case Strong:
		return "Strong"
	case VeryStrong:
		return "Very Strong"
	default:
		return "Unknown"
	}
}

func (s StrengthLevel) Color() string {
	switch s {
	case VeryWeak:
		return "\033[91m" // Red
	case Weak:
		return "\033[91m" // Red
	case Fair:
		return "\033[93m" // Yellow
	case Good:
		return "\033[93m" // Yellow
	case Strong:
		return "\033[92m" // Green
	case VeryStrong:
		return "\033[92m" // Green
	default:
		return "\033[0m" // Reset
	}
}

type PasswordStrength struct {
	Score       int
	Level       StrengthLevel
	Entropy     float64
	Feedback    []string
	TimeToCrack string
}

func AnalyzePasswordStrength(password string) PasswordStrength {
	score := 0
	var feedback []string
	
	length := len(password)
	
	// Length scoring
	if length < 8 {
		feedback = append(feedback, "Use at least 8 characters")
	} else if length < 12 {
		score += 10
		feedback = append(feedback, "Consider using 12+ characters for better security")
	} else if length < 16 {
		score += 20
	} else {
		score += 30
	}
	
	// Character variety scoring
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	
	varietyCount := 0
	if hasLower {
		varietyCount++
		score += 10
	} else {
		feedback = append(feedback, "Add lowercase letters")
	}
	
	if hasUpper {
		varietyCount++
		score += 10
	} else {
		feedback = append(feedback, "Add uppercase letters")
	}
	
	if hasDigit {
		varietyCount++
		score += 10
	} else {
		feedback = append(feedback, "Add numbers")
	}
	
	if hasSymbol {
		varietyCount++
		score += 15
	} else {
		feedback = append(feedback, "Add symbols (!@#$%^&*)")
	}
	
	// Bonus for using all character types
	if varietyCount == 4 {
		score += 10
	}
	
	// Pattern penalties
	if hasRepeatedChars(password) {
		score -= 10
		feedback = append(feedback, "Avoid repeated characters")
	}
	
	if hasSequentialChars(password) {
		score -= 15
		feedback = append(feedback, "Avoid sequential characters (abc, 123)")
	}
	
	if hasCommonPatterns(password) {
		score -= 20
		feedback = append(feedback, "Avoid common patterns")
	}
	
	// Calculate entropy
	entropy := calculateEntropy(password)
	
	// Adjust score based on entropy
	if entropy >= 60 {
		score += 20
	} else if entropy >= 40 {
		score += 10
	} else if entropy < 25 {
		score -= 15
		feedback = append(feedback, "Password is too predictable")
	}
	
	// Ensure score is within bounds
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	
	// Determine strength level
	level := getStrengthLevel(score)
	
	// Generate time to crack estimate
	timeToCrack := estimateTimeToCrack(entropy)
	
	// Add positive feedback for strong passwords
	if score >= 80 && len(feedback) == 0 {
		feedback = append(feedback, "Excellent password strength!")
	}
	
	return PasswordStrength{
		Score:       score,
		Level:       level,
		Entropy:     entropy,
		Feedback:    feedback,
		TimeToCrack: timeToCrack,
	}
}

func calculateEntropy(password string) float64 {
	// Determine character space
	charSpace := 0
	
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		charSpace += 26 // lowercase
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		charSpace += 26 // uppercase
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		charSpace += 10 // digits
	}
	if regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		charSpace += 32 // common symbols
	}
	
	if charSpace == 0 {
		return 0
	}
	
	// Entropy = length * log2(character_space)
	entropy := float64(len(password)) * math.Log2(float64(charSpace))
	
	// Apply penalties for patterns
	if hasRepeatedChars(password) {
		entropy *= 0.8
	}
	if hasSequentialChars(password) {
		entropy *= 0.7
	}
	if hasCommonPatterns(password) {
		entropy *= 0.6
	}
	
	return entropy
}

func hasRepeatedChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}

func hasSequentialChars(password string) bool {
	sequences := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"qwertyuiop", "asdfghjkl", "zxcvbnm",
		"QWERTYUIOP", "ASDFGHJKL", "ZXCVBNM",
	}
	
	lower := strings.ToLower(password)
	for _, seq := range sequences {
		for i := 0; i <= len(seq)-3; i++ {
			if strings.Contains(lower, seq[i:i+3]) {
				return true
			}
		}
	}
	
	// Check for reverse sequences
	for _, seq := range sequences {
		reversed := reverseString(seq)
		for i := 0; i <= len(reversed)-3; i++ {
			if strings.Contains(lower, reversed[i:i+3]) {
				return true
			}
		}
	}
	
	return false
}

func hasCommonPatterns(password string) bool {
	commonPatterns := []string{
		"password", "123456", "qwerty", "admin", "login",
		"welcome", "monkey", "dragon", "master", "shadow",
		"letmein", "football", "iloveyou", "sunshine", "princess",
	}
	
	lower := strings.ToLower(password)
	for _, pattern := range commonPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	
	// Check for simple substitutions
	substitutions := map[string]string{
		"@": "a", "3": "e", "1": "i", "0": "o", "5": "s", "7": "t",
	}
	
	normalized := lower
	for symbol, letter := range substitutions {
		normalized = strings.ReplaceAll(normalized, symbol, letter)
	}
	
	for _, pattern := range commonPatterns {
		if strings.Contains(normalized, pattern) {
			return true
		}
	}
	
	return false
}

func getStrengthLevel(score int) StrengthLevel {
	switch {
	case score < 20:
		return VeryWeak
	case score < 40:
		return Weak
	case score < 60:
		return Fair
	case score < 80:
		return Good
	case score < 95:
		return Strong
	default:
		return VeryStrong
	}
}

func estimateTimeToCrack(entropy float64) string {
	// Assume 1 billion guesses per second (modern hardware)
	guessesPerSecond := 1e9
	
	// Number of possible combinations
	combinations := math.Pow(2, entropy)
	
	// Average time to crack (half the search space)
	secondsToCrack := combinations / (2 * guessesPerSecond)
	
	return formatDuration(secondsToCrack)
}

func formatDuration(seconds float64) string {
	if seconds < 1 {
		return "Instant"
	} else if seconds < 60 {
		return fmt.Sprintf("%.0f seconds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.0f minutes", seconds/60)
	} else if seconds < 86400 {
		return fmt.Sprintf("%.0f hours", seconds/3600)
	} else if seconds < 31536000 {
		return fmt.Sprintf("%.0f days", seconds/86400)
	} else if seconds < 31536000000 {
		return fmt.Sprintf("%.0f years", seconds/31536000)
	} else if seconds < 31536000000000 {
		return fmt.Sprintf("%.0f thousand years", seconds/31536000000)
	} else if seconds < 31536000000000000 {
		return fmt.Sprintf("%.0f million years", seconds/31536000000000)
	} else {
		return fmt.Sprintf("%.0f billion years", seconds/31536000000000000)
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}