package api

import (
	"net/http"

	"github.com/coci/cutme/internal/core/ports"
)

type ShortenerHandler struct {
	Svc ports.Shortener
	Log ports.Logger
}

func NewHandler(svc ports.Shortener, log ports.Logger) *ShortenerHandler {
	return &ShortenerHandler{Svc: svc, Log: log}
}

func (h *ShortenerHandler) ShortLink(w http.ResponseWriter, r *http.Request) {

}

func (h *ShortenerHandler) GetLink(w http.ResponseWriter, r *http.Request) {

}
