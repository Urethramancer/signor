package server

import (
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
	w := web.New(address, port, s.Logger, secure)
	s.webservers[port] = w
	return w
}

// LoadWebServer from JSON config file.
func (s *Server) LoadWebServer(name string) (*web.Web, error) {
	w, err := web.NewFromFile(name, s.Logger)
	if err != nil {
		return nil, err
	}

	s.webservers[w.Port] = w
	return w, err
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
