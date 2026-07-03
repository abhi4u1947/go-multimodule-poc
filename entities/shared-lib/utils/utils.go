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

// Reverse returns s with its runes reversed.
//
// Unreleased: added after v1.1.0 but not yet tagged, used to demonstrate
// Go pseudo-versions when a module is required at a commit with no tag.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
