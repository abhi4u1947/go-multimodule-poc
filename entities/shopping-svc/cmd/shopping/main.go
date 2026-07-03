// Command shopping is the shopping-svc entrypoint. It wires together the
// shared-lib module (logger, config) with shopping-svc's own internal
// service, proving that the root module resolves its cross-module
// dependency correctly.
package main

import (
	"fmt"

	"github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/config"
	"github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/logger"

	"github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/api"
	"github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/internal"
	"github.com/abhi4u1947/go-multimodule-poc/entities/shopping-svc/version"
)

func main() {
	log := logger.New("shopping-svc")

	cfg := config.New()
	cfg.Set("service.name", "shopping-svc")

	svc := internal.NewService(cfg)

	log.Info("starting %s v%s", svc.Name(), version.Version)
	log.Info("api self-check: %s", api.Ping())
	fmt.Println("shopping-svc ready")
}
