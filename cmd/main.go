package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coci/cutme/internal/adapters/api"
	"github.com/coci/cutme/internal/adapters/repositories"
	"github.com/coci/cutme/internal/infra"
	"github.com/coci/cutme/internal/services"
	"github.com/coci/cutme/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	handler := api.ShortenerHandler{
		Svc: services.NewShortenerService(
			repositories.NewLinkRepository(),
			repositories.NewIDGeneratorRepository(cfg),
			cfg,
		),
		Log: infra.NewZapLogger(),
	}

	http.HandleFunc("/short", handler.ShortLink)
	http.HandleFunc("/resolve", handler.GetLink)

	fmt.Println(cfg)
	err = http.ListenAndServe(cfg.BaseURL, nil)

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
