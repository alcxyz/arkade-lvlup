package tools

import (
	"strings"
)

// ContainsElement checks if a slice contains a given element.
func ContainsElement(slice []string, elem string) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}

// PopulateArray splits a comma-separated string into a slice of strings.
func PopulateArray(input string) []string {
	return strings.Split(strings.TrimSpace(input), ",")
}
