package server

import (
	"sync"

	"github.com/Urethramancer/signor/log"
)

// Server structure for a web server.
type Server struct {
	sync.RWMutex
	sync.WaitGroup
	// Name for logging purposes and management.
	Name string
	web  map[string]*Site
	// Very lazy default logging
	l    func(string, ...interface{})
	e    func(string, ...interface{})
	quit chan bool
}

// New server instance creation.
func New(name string) *Server {
	s := Server{
		Name: name,
		quit: make(chan bool, 1),
		web:  make(map[string]*Site),
		l:    log.Default.TMsg,
		e:    log.Default.TErr,
	}

	return &s
}

// Start all configured sub-servers.
func (s *Server) Start() error {
	s.l("Starting server '%s'.", s.Name)

	s.Add(1)
	go func() {
		<-s.quit
		s.l("Quitting server '%s'.", s.Name)
		s.Done()
	}()
	return nil
}

// Stop all sub-servers.
func (s *Server) Stop() {
	s.quit <- true
	s.Wait()
}
