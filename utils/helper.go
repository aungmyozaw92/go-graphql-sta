package utils

import (
	"unicode"
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

// returns slice removing duplicate elements
func UniqueSlice[T comparable](slice []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, elm := range slice {
		if _, ok := inResult[elm]; !ok {
			// if not exists in map, append it, otherwise do nothing
			inResult[elm] = true
			result = append(result, elm)
		}
	}
	return result
}

var CountryCode = "MM"
