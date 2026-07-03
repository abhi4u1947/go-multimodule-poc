// Package utils holds small shared helpers with no external dependencies.
package utils

import "strings"

// Slugify lowercases s and replaces whitespace runs with a single hyphen.
func Slugify(s string) string {
	fields := strings.Fields(s)
	return strings.ToLower(strings.Join(fields, "-"))
}

// Truncate shortens s to at most n runes, appending "..." if it was cut.
func Truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
