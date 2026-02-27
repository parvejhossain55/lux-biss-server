package validators

import (
	"regexp"
	"unicode"
)

// PasswordStrength represents the strength level of a password
type PasswordStrength int

const (
	Weak PasswordStrength = iota
	Medium
	Strong
)

// PasswordValidationResult contains the result of password validation
type PasswordValidationResult struct {
	IsValid  bool
	Strength PasswordStrength
	Errors   []string
	Score    int
}

// ValidatePassword validates password strength and returns detailed results
func ValidatePassword(password string) *PasswordValidationResult {
	result := &PasswordValidationResult{
		IsValid: true,
		Errors:  []string{},
		Score:   0,
	}

	// Length validation
	if len(password) < 8 {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must be at least 8 characters long")
	} else if len(password) >= 8 {
		result.Score += 1
	}

	if len(password) >= 12 {
		result.Score += 1
	}

	// Character type validation
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasLower {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one lowercase letter")
	} else {
		result.Score += 1
	}

	if !hasUpper {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one uppercase letter")
	} else {
		result.Score += 1
	}

	if !hasNumber {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one number")
	} else {
		result.Score += 1
	}

	if !hasSpecial {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password must contain at least one special character")
	} else {
		result.Score += 1
	}

	// Common patterns check
	if isCommonPattern(password) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password contains common patterns and is not secure")
		result.Score -= 2
	}

	// Sequential characters check
	if hasSequentialChars(password) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Password contains sequential characters")
		result.Score -= 1
	}

	// Determine strength
	if result.Score >= 6 {
		result.Strength = Strong
	} else if result.Score >= 4 {
		result.Strength = Medium
	} else {
		result.Strength = Weak
	}

	return result
}

// isCommonPattern checks for common password patterns
func isCommonPattern(password string) bool {
	commonPatterns := []string{
		"password", "123456", "qwerty", "admin", "letmein",
		"welcome", "monkey", "dragon", "master", "hello",
		"freedom", "whatever", "qazwsx", "trustno1", "123qwe",
		"1q2w3e4r", "abc123", "password123",
	}

	lowerPassword := toLower(password)
	for _, pattern := range commonPatterns {
		if contains(lowerPassword, pattern) {
			return true
		}
	}

	return false
}

// hasSequentialChars checks for sequential characters
func hasSequentialChars(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		// Check for sequential numbers
		if password[i] >= '0' && password[i] <= '7' {
			if password[i+1] == password[i]+1 && password[i+2] == password[i]+2 {
				return true
			}
		}
		// Check for sequential letters
		if password[i] >= 'a' && password[i] <= 'y' {
			if password[i+1] == password[i]+1 && password[i+2] == password[i]+2 {
				return true
			}
		}
		if password[i] >= 'A' && password[i] <= 'Y' {
			if password[i+1] == password[i]+1 && password[i+2] == password[i]+2 {
				return true
			}
		}
	}

	return false
}

// Helper functions
func toLower(s string) string {
	result := make([]rune, len([]rune(s)))
	for i, r := range []rune(s) {
		result[i] = unicode.ToLower(r)
	}
	return string(result)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

// GetPasswordRequirements returns a list of password requirements
func GetPasswordRequirements() []string {
	return []string{
		"At least 8 characters long",
		"Contains at least one lowercase letter",
		"Contains at least one uppercase letter",
		"Contains at least one number",
		"Contains at least one special character",
		"Does not contain common patterns",
		"Does not contain sequential characters",
	}
}

// GetStrengthDescription returns a human-readable description of password strength
func GetStrengthDescription(strength PasswordStrength) string {
	switch strength {
	case Strong:
		return "Strong"
	case Medium:
		return "Medium"
	case Weak:
		return "Weak"
	default:
		return "Unknown"
	}
}
