package utils

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/ttacon/libphonenumber"
)

func NewTrue() *bool {
	b := true
	return &b
}

func NewFalse() *bool {
	b := false
	return &b
}

// turn salesInvoice to SalesInvoice
func UppercaseFirst(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// turn ToggleActive to toggleActive
func LowercaseFirst(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

var CountryCode = "MM"

func IsValidEmail(email string) bool {
	// Basic email validation regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}


func ValidatePhoneNumber(phoneNumber, countryCode string) error {
	p, err := libphonenumber.Parse(phoneNumber, countryCode)
	if err != nil {
		return err // Phone number is invalid
	}

	if !libphonenumber.IsValidNumber(p) {
		return fmt.Errorf("phone number is not valid")
	}

	return nil // Phone number is valid for the specified country code
}