package main

import (
	"log"

	"github.com/ivanvillarroel/go_jwt_gin/internal/app"
	"github.com/ivanvillarroel/go_jwt_gin/internal/config"
)

func main() {
	cfg := config.Load()

	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
