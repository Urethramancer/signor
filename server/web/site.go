package web

import (
	"fmt"
	"net/http"

	"github.com/Urethramancer/signor/log"
)

// Site is the configuration for a specific domain.
type Site struct {
	// Domain to respond to.
	Domain string `json:"domain"`
	// Certificate PEM file to load. Absolute path.
	Certificate string `json:"certificate,omitempty"`
	// Key PEM file to load. Absolute path.
	Key string `json:"key,omitempty"`
	// Owner is a user in the database.
	Owner string `json:"owner"`

	//
	// Internals
	//
	log.LogShortcuts
}

// SetLogger changes the logger object and sets the message shortcuts for convenience.
func (s *Site) SetLogger(l *log.Logger) {
	s.Logger = l
	s.L = l.TMsg
	s.E = l.TErr
}

func enableCors(wr *http.ResponseWriter) {
	(*wr).Header().Set("Access-Control-Allow-Origin", "*")
}

// DefaultHandler for testing.
func (s *Site) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s: It sort of works!", r.Host)
}

func (s *Site) favicon(w http.ResponseWriter, r *http.Request) error {
	s.L("Favicon for site %s", s.Domain)
	return nil
}
