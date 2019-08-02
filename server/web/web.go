package web

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Urethramancer/signor/files"
	"github.com/Urethramancer/signor/log"
)

// Web server configuration and runtime data.
type Web struct {
	sync.RWMutex
	sync.WaitGroup
	log.LogShortcuts
	http.Server

	// Address to bind the web server to.
	Address string `json:"address"`
	// Port is optional, and will default to 80 or 443, depending on presence of certificates.
	Port string `json:"port"`
	// SitePath is the directory where site configurations are stored.
	SitePath string `json:"sites"`
	// Secure servers prepare a default TLS configuration.
	Secure bool `json:"secure"`
	// Timeouts for connections and shutdown.
	Timeouts Timeouts `json:"timeouts"`

	//
	// Internals
	//
	sites   map[string]*Site
	running bool
}

// Timeouts for web server.
type Timeouts struct {
	// Idle timeout defaults to 30 seconds.
	Idle time.Duration `json:"idle"`
	// Read timeout defaults to 10 seconds.
	Read time.Duration `json:"read"`
	// Write timeout defaults to 10 seconds.
	Write time.Duration `json:"write"`
	// Shutdown timeout defaults to 3 seconds.
	Shutdown time.Duration `json:"shutdown"`
}

// New creates a web server, configured with reasonable timeouts and a default handlers.
// Use AddSite() or LoadSites() to add handlers for domains.
func New(address, port string, logger *log.Logger, secure bool) *Web {
	w := Web{
		Address: address,
		Port:    port,
		Secure:  secure,
	}
	w.Init(logger)
	return &w
}

// Init sets up basics for a web server.
func (w *Web) Init(l *log.Logger) {
	if w.Port == "" {
		if !w.Secure {
			w.Port = "80"
		} else {
			w.Port = "443"
		}
	}

	w.sites = make(map[string]*Site)
	w.Logger = l
	w.L = l.TMsg
	w.E = l.TErr

	w.IdleTimeout = time.Second * 30
	w.ReadTimeout = time.Second * 10
	w.WriteTimeout = time.Second * 10

	if w.Secure {
		w.TLSConfig = &tls.Config{
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
	}
	http.HandleFunc("/", w.defaultHandler)
}

// NewFromFile creates a web server based on a JSON configuration file.
func NewFromFile(name string, logger *log.Logger) (*Web, error) {
	in, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var w Web
	err = json.Unmarshal(in, &w)
	if err != nil {
		return nil, err
	}

	w.IdleTimeout = nonZeroDuration(w.Timeouts.Idle, 30)
	w.ReadTimeout = nonZeroDuration(w.Timeouts.Read, 10)
	w.WriteTimeout = nonZeroDuration(w.Timeouts.Write, 10)

	w.Init(logger)
	return &w, nil
}

// SetLogger changes the logger object and sets the message shortcuts for convenience.
func (w *Web) SetLogger(l *log.Logger) {
	w.Lock()
	defer w.Unlock()
	w.Logger = l
	w.L = l.TMsg
	w.E = l.TErr
}

// AddCertificate from loaded certificate.
func (w *Web) AddCertificate(cert tls.Certificate) {
	w.TLSConfig.Certificates = append(w.TLSConfig.Certificates, cert)
	w.TLSConfig.BuildNameToCertificate()
}

// RebuildCertificates reloads the certificates from all sites.
func (w *Web) RebuildCertificates() error {
	for _, s := range w.sites {
		cert, err := tls.LoadX509KeyPair(s.Certificate, s.Key)
		if err != nil {
			return err
		}

		w.AddCertificate(cert)
	}
	w.TLSConfig.BuildNameToCertificate()
	return nil
}

// AddSite to a web server.
// This is done on the fly without need of any restarting.
func (w *Web) AddSite(s *Site) error {
	_, ok := w.sites[s.Domain]
	if ok {
		return errors.New(ErrSiteExists)
	}

	if w.Secure {
		cert, err := tls.LoadX509KeyPair(s.Certificate, s.Key)
		if err != nil {
			return err
		}

		//TODO: Watch key+cert for changes and automatically reload.
		//TODO: Let's Encrypt support.
		w.AddCertificate(cert)
	}
	w.sites[s.Domain] = s
	http.HandleFunc(s.Domain+"/", s.DefaultHandler)
	return nil
}

// Start the webserver.
func (w *Web) Start() {
	if w.running {
		return
	}

	w.Lock()
	defer w.Unlock()
	w.running = true

	var err error
	addr := net.JoinHostPort(w.Address, w.Port)
	w.L("Starting web server on %s (secure=%t)", addr, w.Secure)
	if w.Secure {
		err = w.RebuildCertificates()
		if err != nil {
			w.E("Certificate error: %s", err.Error())
			return
		}

		listener, err := tls.Listen("tcp", addr, w.TLSConfig)
		if err != nil {
			w.running = false
			w.E("TLS listener error: %s", err.Error())
			return
		}

		w.Add(1)
		err = w.Serve(listener)
		if err != nil {
			w.E("Web server error: %s", err.Error())
		}
		w.Done()
		return
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		w.running = false
		w.E("Listener error: %s", err.Error())
		return
	}

	w.Add(1)
	err = w.Serve(listener)
	if err != nil {
		w.E("Web server error: %s", err.Error())
	}
	w.Done()
}

// Stop the webserver and try to wait until all connections are done.
// Give up after half a configured timeout.
func (w *Web) Stop() error {
	w.L("Stopping web server on %s:%s", w.Address, w.Port)
	ctx, cancel := context.WithTimeout(context.Background(), nonZeroDuration(w.Timeouts.Shutdown, time.Second*30))
	defer cancel()
	err := w.Shutdown(ctx)
	if err != nil {
		return err
	}

	w.running = false
	return nil
}

// defaultHandler when no sites are configured on a requested URL.
func (w *Web) defaultHandler(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(wr, "No site configured on %s", r.Host)
}

// LoadSites from JSON files in the site config directory (non-recursively).
func (w *Web) LoadSites() error {
	if !files.DirExists(w.SitePath) {
		return os.ErrNotExist
	}

	dir, err := ioutil.ReadDir(w.SitePath)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		fn := filepath.Join(w.SitePath, fi.Name())
		in, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		var site Site
		err = json.Unmarshal(in, &site)
		if err != nil {
			return err
		}

		w.L("Web: Loaded %s", site.Domain)
		err = w.AddSite(&site)
		if err != nil {
			return err
		}
	}
	return nil
}
