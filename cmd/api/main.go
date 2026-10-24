package main

import (
	"log"

	"github.com/sahilrana7582/user-svc/internal/config"
)

func main() {
	cfg, err := config.LoadViper()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	log.Printf("starting %s on %s:%s", cfg.App.Name, cfg.App.Host, cfg.App.Port)
}
