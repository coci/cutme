package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coci/cutme/internal/adapters/api"
	"github.com/coci/cutme/internal/adapters/repositories"
	"github.com/coci/cutme/internal/infra/config"
	"github.com/coci/cutme/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	handler := api.ShortenerHandler{
		Svc: services.NewShortenerService(
			repositories.NewLinkRepository(cfg),
			repositories.NewIDGeneratorRepository(cfg),
			cfg,
		),
		Log: services.NewZapLogger(),
	}

	http.HandleFunc("/short", handler.ShortLink)
	http.HandleFunc("/resolve", handler.GetLink)

	fmt.Printf("Starting server on: %s\n", cfg.BaseURL)

	err = http.ListenAndServe(cfg.BaseURL, nil)

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
