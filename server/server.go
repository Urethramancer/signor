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
	log.LogShortcuts

	// Name for logging purposes and management.
	Name       string
	webservers map[string]*web.Web
	quit       chan bool
}

// New server instance creation.
func New(name string) *Server {
	s := Server{
		Name:       name,
		webservers: make(map[string]*web.Web),
		quit:       make(chan bool, 1),
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
func (s *Server) Stop() error {
	s.quit <- true
	err := s.StopWeb()
	if err != nil {
		return err
	}
	s.Wait()
	return nil
}

// SetLogger changes the logger object and sets the message shortcuts for convenience.
func (s *Server) SetLogger(l *log.Logger) {
	s.Logger = l
	s.L = l.TMsg
	s.E = l.TErr
	for _, w := range s.webservers {
		w.SetLogger(l)
	}
}

// AddWebServer to the server.
func (s *Server) AddWebServer(address, port string, secure bool) *web.Web {
	var w *web.Web
	if secure {
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
		w = web.New(address, port, s.Logger, cfg)
	} else {
		w = web.New(address, port, s.Logger, nil)
	}

	s.webservers[port] = w
	return w
}

// StartWeb starts all web servers.
func (s *Server) StartWeb() {
	for _, w := range s.webservers {
		go w.Start()
	}
}

// StopWeb stops all web servers.
func (s *Server) StopWeb() error {
	var err error
	for _, w := range s.webservers {
		err = w.Stop()
		if err != nil {
			return err
		}
		w.Wait()
	}
	return nil
}
