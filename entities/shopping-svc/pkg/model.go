// Package pkg holds shopping-svc types that are safe for external modules
// to import (unlike internal/).
package pkg

// Item is a minimal shopping catalog entry.
type Item struct {
	SKU   string
	Name  string
	Price float64
}
