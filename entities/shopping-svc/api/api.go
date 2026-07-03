// Package api holds shopping-svc's HTTP handler surface.
//
// Kept intentionally minimal: this PoC is about module resolution, not
// building a real shopping API.
package api

// Ping is a placeholder handler used by cmd/shopping to prove the api
// package participates in the shopping-svc module's build graph.
func Ping() string {
	return "pong"
}
