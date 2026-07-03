// Package internal holds shopping-svc business logic that must never be
// imported outside the shopping-svc module tree (enforced by Go's internal/
// visibility rule).
package internal

import (
	"github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/config"
	"github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/utils"
)

// Service is the shopping-svc core service, wired with shared-lib config.
type Service struct {
	cfg *config.Config
}

// NewService returns a Service backed by cfg.
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// Name returns the configured service name, defaulting to "shopping-svc".
func (s *Service) Name() string {
	return s.cfg.GetOrDefault("service.name", "shopping-svc")
}

// IsShoppingRelated reports whether text mentions this service by name.
//
// Added in shopping-svc v1.1.0 to exercise shared-lib's new
// utils.Contains helper (shared-lib v1.1.0).
func (s *Service) IsShoppingRelated(text string) bool {
	return utils.Contains(text, s.Name())
}
