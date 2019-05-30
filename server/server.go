package server

import (
	"sync"

	"github.com/Urethramancer/slog"
)

type Server struct {
	sync.RWMutex
	sync.WaitGroup
	// Name for logging purposes and management.
	Name string
	quit chan bool
	web  map[string]*Site
}

// New creates a server instance.
func New(name string) *Server {
	s := Server{
		Name: name,
		quit: make(chan bool, 1),
		web:  make(map[string]*Site),
	}

	return &s
}

// Start all configured sub-servers.
func (s *Server) Start() error {
	slog.TMsg("Starting server '%s'.", s.Name)

	s.Add(1)
	go func() {
		<-s.quit
		slog.TMsg("Quitting server '%s'.", s.Name)
		s.Done()
	}()
	return nil
}

// Stop all sub-servers.
func (s *Server) Stop() {
	s.quit <- true
	s.Wait()
}
