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

// Contains reports whether needle appears in haystack, case-insensitively.
//
// Added in shared-lib v1.1.0 to demonstrate a backward-compatible minor
// release that downstream modules can pick up via `go get -u`.
func Contains(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}
