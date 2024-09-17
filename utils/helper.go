package utils

import (
	"fmt"
	"regexp"

	"github.com/ttacon/libphonenumber"
)


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