package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Urethramancer/signor/log"
)

type Web struct {
	// Address to bind the web server to.
	Address string `json:"address"`
	// SecurePort is optional, and will default to 443.
	SecurePort string `json:"secureport,omitempty"`
	// OtherPort is optional, and will default to 80.
	OtherPort string `json:"otherport,omitempty"`

	//
	// Internals
	//
	secureSites map[string]*Site
	sites       map[string]*Site
	server      *http.Server
	log.LogShortcuts
}

func New(address, secureport, otherport string, logger *log.Logger, cfg *tls.Config) *Web {
	w := Web{
		sites:      make(map[string]*Site),
		Address:    address,
		SecurePort: secureport,
		OtherPort:  otherport,
	}
	w.Logger = logger
	w.L = logger.TMsg
	w.E = logger.TErr

	w.server = &http.Server{
		IdleTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		TLSConfig:    cfg,
	}

	http.HandleFunc("/", w.defaultHandler)
	return &w
}

// SetLogger changes the logger object and sets the message shortcuts for convenience.
func (w *Web) SetLogger(l *log.Logger) {
	w.Logger = l
	w.L = l.TMsg
	w.E = l.TErr
}

// AddCertificate from loaded certificate.
func (w *Web) AddCertificate(cert tls.Certificate) {
	w.server.TLSConfig.Certificates = append(w.server.TLSConfig.Certificates, cert)
	w.server.TLSConfig.BuildNameToCertificate()
}

// RebuildCertificates reloads the certificates from all sites.
func (w *Web) RebuildCertificates() {
	w.server.TLSConfig.BuildNameToCertificate()
}

// AddSite to a web server.
// This is done on the fly without need of any restarting.
func (w *Web) AddSite(s *Site) error {
	cert, err := tls.LoadX509KeyPair(s.Certificate, s.Key)
	if err != nil {
		return err
	}

	//TODO: Watch key+cert for changes and automatically reload.
	//TODO: Let's Encrypt support.
	w.AddCertificate(cert)
	w.sites[s.Domain] = s
	http.HandleFunc(s.Domain+"/", s.DefaultHandler)
	return nil
}

// Start the webserver.
func (w *Web) Start() error {
	w.RebuildCertificates()
	var err error
	listener, err := tls.Listen("tcp", net.JoinHostPort(w.Address, w.Port), w.server.TLSConfig)
	if err != nil {
		return err
	}

	return w.server.Serve(listener)
}

// Stop the webserver and wait till all connections are done.
func (w *Web) Stop() error {
	return w.server.Shutdown(context.Background())
}

// defaultHandler when no sites are configured on a requested URL.
func (web *Web) defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "No site configured on %s", r.Host)
}
