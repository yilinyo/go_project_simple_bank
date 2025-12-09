package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile("^[a-zA-Z0-9_]+$").MatchString
	isValidFullName = regexp.MustCompile("^[a-zA-Z\\s ]+$").MatchString
)

func ValidateString(value string, minLen int, maxLen int) error {
	n := len(value)
	if n < minLen {
		return fmt.Errorf("string length %d is less than minimum length %d", n, minLen)
	}
	if n > maxLen {
		return fmt.Errorf("string length %d exceeds maximum length %d", n, maxLen)
	}
	return nil
}

func ValidateUsername(value string) error {

	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("username %s contains invalid characters", value)
	}
	return nil
}

func ValidatePassword(value string) error {

	return ValidateString(value, 6, 32)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 5, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email format: %v", err)
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("full name %s contains invalid characters", value)
	}
	return nil
}
