package server

import (
	"net/http"

	"github.com/Urethramancer/slog"
	"github.com/gorilla/mux"
)

// Site holds one domain's key configuration.
type Site struct {
	mux.Router
	// Domain to respond to.
	Domain string `json:"domain"`
	// Owner is a user in the database.
	Owner string `json:"owner"`
	// Certificate PEM file to load. Absolute path.
	Certificate string `json:"certificate,omitempty"`
	// Key PEM file to load. Absolute path.
	Key string `json:"key,omitempty"`

	// Internal data for the instance.
	url string
}

func (site *Site) favicon(w http.ResponseWriter, r *http.Request) error {
	slog.Msg("Favicon for site %s", site.Domain)
	return nil
}
