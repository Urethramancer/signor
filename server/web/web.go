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

type Web struct {
	sync.RWMutex
	sync.WaitGroup
	// Address to bind the web server to.
	Address string `json:"address"`
	// Port is optional, and will default to 80 or 443, depending on presence of certificates.
	Port string `json:"secureport,omitempty"`

	//
	// Internals
	//
	secureSites map[string]*Site
	sites       map[string]*Site
	server      *http.Server
	log.LogShortcuts
	secure  bool
	running bool
}

// New creates a web server, configured with reasonable timeouts and a default handlers.
// Use AddSite() to add handlers for domains.
func New(address, port string, logger *log.Logger, cfg *tls.Config) *Web {
	if port == "" {
		if cfg == nil {
			port = "80"
		} else {
			port = "443"
		}
	}

	w := Web{
		sites:   make(map[string]*Site),
		Address: address,
		Port:    port,
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

	if cfg != nil {
		w.secure = true
	}

	http.HandleFunc("/", w.defaultHandler)
	return &w
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
	w.Lock()
	defer w.Unlock()
	w.server.TLSConfig.Certificates = append(w.server.TLSConfig.Certificates, cert)
	w.server.TLSConfig.BuildNameToCertificate()
}

// RebuildCertificates reloads the certificates from all sites.
func (w *Web) RebuildCertificates() error {
	w.Lock()
	defer w.Unlock()
	for _, s := range w.secureSites {
		cert, err := tls.LoadX509KeyPair(s.Certificate, s.Key)
		if err != nil {
			return err
		}

		w.AddCertificate(cert)
	}
	w.server.TLSConfig.BuildNameToCertificate()
	return nil
}

// AddSite to a web server.
// This is done on the fly without need of any restarting.
func (w *Web) AddSite(s *Site) error {
	_, ok := w.sites[s.Domain]
	if ok {
		return errors.New(ErrSiteExists)
	}

	if w.secure {
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
func (w *Web) Start() error {
	if w.running {
		return nil
	}

	w.Lock()
	defer w.Unlock()
	w.running = true

	var err error
	addr := net.JoinHostPort(w.Address, w.Port)
	w.L("Starting web server on %s", addr)
	if w.secure {
		w.RebuildCertificates()
		listener, err := tls.Listen("tcp", addr, w.server.TLSConfig)
		if err != nil {
			w.running = false
			return err
		}

		w.Add(1)
		return w.server.Serve(listener)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		w.running = false
		return err
	}

	w.Add(1)
	return w.server.Serve(listener)
}

// Stop the webserver and try to wait till all connections are done.
// Give up after half a second.
func (w *Web) Stop() error {
	w.L("Stopping web server on %s:%s", w.Address, w.Port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	err := w.server.Shutdown(ctx)
	w.Done()
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

// LoadSites loads sites from JSON files in the specified path (non-recursively).
func (w *Web) LoadSites(path string) error {
	if !files.DirExists(path) {
		return os.ErrNotExist
	}

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		fn := filepath.Join(path, fi.Name())
		in, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		var site Site
		json.Unmarshal(in, &site)
		w.L("Web: Loaded %s", site.Domain)
		w.AddSite(&site)
	}
	return nil
}
