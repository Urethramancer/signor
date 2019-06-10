package main

import (
	"os"

	"github.com/Urethramancer/daemon"
	"github.com/Urethramancer/signor/server"
)

func main() {
	var err error
	s := server.New("example")
	err = s.Start()
	if err != nil {
		s.E("Error starting server: %s", err.Error())
		os.Exit(2)
	}

	// Create a plain web server with two sites.
	w := s.AddWebServer("127.0.0.1", "11000", false)
	err = w.LoadSites("sites")
	if err != nil {
		s.L("Error loading sites: %s", err.Error())
		os.Exit(2)
	}

	s.StartWeb()
	if err != nil {
		s.E("Error starting web server: %s", err.Error())
		os.Exit(2)
	}

	// Wait for ctrl-c
	<-daemon.BreakChannel()
	err = s.Stop()
	if err != nil {
		s.E("Error stopping server(s): %s", err.Error())
		os.Exit(2)
	}
	s.L("Stopped.")
}
