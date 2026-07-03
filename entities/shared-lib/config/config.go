// Package config provides a tiny in-memory key/value config store shared by
// every service module in the monorepo.
package config

// Config is a trivial string-keyed configuration map.
type Config struct {
	values map[string]string
}

// New returns an empty Config.
func New() *Config {
	return &Config{values: make(map[string]string)}
}

// Set stores a value under key.
func (c *Config) Set(key, value string) {
	c.values[key] = value
}

// Get returns the value stored under key, and whether it was present.
func (c *Config) Get(key string) (string, bool) {
	v, ok := c.values[key]
	return v, ok
}

// GetOrDefault returns the stored value, or def if key is absent.
func (c *Config) GetOrDefault(key, def string) string {
	if v, ok := c.values[key]; ok {
		return v
	}
	return def
}
