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
	w := s.AddWebServer("127.0.0.1", "10000")
	site := &web.Site{
		Domain:      "localhost",
		Certificate: "cert.pem",
		Key:         "key.pem",
		Owner:       "orb",
	}

	var err error
	err = w.AddSite(site)
	if err != nil {
		s.E("Error adding web domain: %s", err.Error())
		os.Exit(2)
	}

	site = &web.Site{
		Domain:      "localhost.com",
		Certificate: "cert.pem",
		Key:         "key.pem",
		Owner:       "orb",
	}
	err = w.AddSite(site)
	if err != nil {
		s.E("Error adding web domain: %s", err.Error())
		os.Exit(2)
	}

	err = w.Start()
	if err != nil {
		s.E("Error starting web server: %s", err.Error())
		os.Exit(2)
	}

	<-daemon.BreakChannel()
	err = w.Stop()
	if err != nil {
		s.E("Error stopping web server: %s", err.Error())
		os.Exit(2)
	}

	s.Stop()
}
