// Package main is the entry point for the application server.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/fox-gonic/fox"

	"github.com/miclle/go-single-page-binary-app/internal/config"
	"github.com/miclle/go-single-page-binary-app/internal/handler"
	"github.com/miclle/go-single-page-binary-app/internal/service"
)

var (
	// CommitID is the git commit hash, injected at build time via ldflags.
	CommitID = "dev"
	// BuildTime is the build timestamp, injected at build time via ldflags.
	BuildTime = ""
)

func main() {
	configPath := flag.String("c", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()
	svc, err := service.New(ctx, cfg.Driver, cfg.DSN)
	if err != nil {
		log.Fatalf("init service: %v", err)
	}

	engine := fox.Default()
	ctrl := handler.New(svc)
	ctrl.RegisterRoutes(engine)

	log.Printf("server starting on %s (commit=%s, built=%s)", cfg.Addr, CommitID, BuildTime)
	if err := engine.Run(cfg.Addr); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
