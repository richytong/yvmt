package main

import (
	"regexp"
	"strings"
)

var alphaNumericRegex = regexp.MustCompile("[^a-zA-Z0-9\\s]+")

// "here's a string -> []slice{heres, a, string}"
func extractLowerAlphanumericFields(s string) []string {
	return strings.Fields(strings.ToLower(alphaNumericRegex.ReplaceAllString(s, "")))
}
