package utils

import (
	"unicode"
)

// ContainsUppercase checks if the string contains at least one uppercase letter.
func ContainsUppercase(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

// ContainsLowercase checks if the string contains at least one lowercase letter.
func ContainsLowercase(s string) bool {
	for _, char := range s {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

// ContainsNumber checks if the string contains at least one numeric character.
func ContainsNumber(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
