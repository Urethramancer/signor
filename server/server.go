package server

import (
	"crypto/tls"
	"sync"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/server/web"
)

// Server structure for a web server.
type Server struct {
	sync.RWMutex
	sync.WaitGroup
	// Name for logging purposes and management.
	Name string
	web  *web.Web
	log.LogShortcuts
	quit chan bool
}

// New server instance creation.
func New(name string) *Server {
	s := Server{
		Name: name,
		quit: make(chan bool, 1),
	}
	s.Logger = log.Default
	s.L = log.Default.TMsg
	s.E = log.Default.TErr

	return &s
}

// Start all configured sub-servers.
func (s *Server) Start() error {
	s.L("Starting server '%s'.", s.Name)

	s.Add(1)
	go func() {
		<-s.quit
		s.L("Quitting server '%s'.", s.Name)
		s.Done()
	}()
	return nil
}

// Stop all sub-servers.
func (s *Server) Stop() {
	s.quit <- true
	s.Wait()
}

// SetLogger changes the logger object and sets the message shortcuts for convenience.
func (s *Server) SetLogger(l *log.Logger) {
	s.Logger = l
	s.L = l.TMsg
	s.E = l.TErr
	if s.web != nil {
		s.web.SetLogger(l)
	}
}

// AddWebServer to the server.
func (s *Server) AddWebServer(address, secureport, otherport string) *web.Web {
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		NextProtos:               []string{"http/1.1"},
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	w := web.New(address, secureport, otherport, s.Logger, cfg)
	s.web = w
	return w
}
