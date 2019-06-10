package main

import (
	"os"

	"github.com/Urethramancer/daemon"
	"github.com/Urethramancer/signor/server"
	"github.com/Urethramancer/signor/server/web"
)

func main() {
	s := server.New("example")
	s.Start()
	// Create a plain web server with two sites.
	w := s.AddWebServer("127.0.0.1", "11000", false)
	site := &web.Site{Domain: "localhost"}

	var err error
	err = w.AddSite(site)
	if err != nil {
		s.E("Error adding web domain: %s", err.Error())
		os.Exit(2)
	}

	site = &web.Site{Domain: "localhost.com"}
	err = w.AddSite(site)
	if err != nil {
		s.E("Error adding web domain: %s", err.Error())
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
