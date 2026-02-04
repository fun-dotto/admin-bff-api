package handler

import (
	"net/http"
	"time"

	api "github.com/fun-dotto/api-template/generated"
)

// Handler implements api.ServerInterface
type Handler struct {
	announcementAPIURL string
	httpClient         *http.Client
}

// NewHandler creates a new Handler instance
func NewHandler(announcementAPIURL string) *Handler {
	return &Handler{
		announcementAPIURL: announcementAPIURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Ensure Handler implements api.ServerInterface
var _ api.ServerInterface = (*Handler)(nil)
