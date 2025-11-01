package api

import (
	"encoding/json"
	"net/http"

	"github.com/coci/cutme/internal/core/domain"
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
	var shortLink domain.Link

	err := json.NewDecoder(r.Body).Decode(&shortLink)
	if err != nil {
		h.Log.Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortLink.Code, _ = h.Svc.Shorten(shortLink.Link)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shortLink); err != nil {
		h.Log.Error("failed to encode response")
	}
}

func (h *ShortenerHandler) GetLink(w http.ResponseWriter, r *http.Request) {

}
