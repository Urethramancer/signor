package main

import (
	"os"

	"github.com/Urethramancer/daemon"
	"github.com/Urethramancer/signor/server"
)

func main() {
	var err error
	s := server.New("secure example")
	err = s.Start()
	if err != nil {
		s.E("Error starting server: %s", err.Error())
		os.Exit(2)
	}

	// Create a plain web server with two sites.
	w, err := s.LoadWebServer("web.json")
	if err != nil {
		s.L("Error starting web server: %s", err.Error())
		os.Exit(2)
	}

	err = w.LoadSites()
	if err != nil {
		s.L("Error loading sites: %s", err.Error())
		os.Exit(2)
	}

	s.StartWeb()
	// Wait for ctrl-c
	<-daemon.BreakChannel()
	err = s.Stop()
	if err != nil {
		s.E("Error stopping server(s): %s", err.Error())
		os.Exit(2)
	}
	s.L("Stopped.")
}
