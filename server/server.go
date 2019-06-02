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
	Name   string
	web    map[string]*Site
	logger *log.Logger
	// L is the regular message output command from the Logger object.
	L func(string, ...interface{})
	// E us the error message output command from the Logger object.
	E    func(string, ...interface{})
	quit chan bool
}

// New server instance creation.
func New(name string) *Server {
	s := Server{
		Name:   name,
		quit:   make(chan bool, 1),
		web:    make(map[string]*Site),
		logger: log.Default,
		L:      log.Default.TMsg,
		E:      log.Default.TErr,
	}

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
	s.logger = l
	s.L = l.TMsg
	s.E = l.TErr
}
