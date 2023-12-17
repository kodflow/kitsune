package http

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/core/certs"
	"github.com/kodflow/kitsune/src/internal/core/server/api"
	"github.com/kodflow/kitsune/src/internal/core/server/handler"
	"github.com/kodflow/kitsune/src/internal/kernel/errors"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"golang.org/x/net/http2"
)

// Server represents an HTTP server that handles both standard and secure connections.
// It encapsulates the functionality of two engines, one for handling standard HTTP
// connections and the other for secure HTTPS connections, along with a router for API routing.
type Server struct {
	standard *Engine     // The engine for handling standard HTTP connections.
	secure   *Engine     // The engine for handling secure HTTPS connections.
	router   *api.Router // The router for managing API endpoints.
}

// ServerCfg holds configuration data for the HTTP server.
// This structure includes details such as the server's domain, subdomains, and port numbers
// for both HTTP and HTTPS connections.
type ServerCfg struct {
	DOMAIN string   // The domain of the server.
	SUBS   []string // The subdomains of the server.
	HTTP   int      // The port number for HTTP connections.
	HTTPS  int      // The port number for HTTPS connections.
}

// Engine represents an HTTP engine.
// It is responsible for managing network listeners and the HTTP server for either standard
// or secure connections, keeping track of its running status and configuration details
// like port, domain, and subdomains.
type Engine struct {
	PORT     int          // The port number for the engine.
	DOMAIN   string       // The domain of the engine.
	SUBS     []string     // The subdomains of the engine.
	listener net.Listener // The network listener for the engine.
	server   *http.Server // The HTTP server instance.
	running  bool         // Indicates if the engine is currently running.
}

// newServerConfig creates a new HTTP server configuration.
// It sets up various timeouts and limits for the server. If a TLS configuration is provided,
// it is applied to enable HTTPS.
//
// Parameters:
// - tls: *tls.Config Optional TLS configuration for HTTPS.
//
// Returns:
// - *http.Server: Configured HTTP server.
func newServerConfig(tls *tls.Config) *http.Server {
	srv := &http.Server{
		Handler:           http.HandlerFunc(handler.HTTPHandler),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	if tls != nil {
		srv.TLSConfig = tls
	}

	return srv
}

// NewServer creates a new HTTP server based on the provided configuration.
// It initializes a standard and, if specified, a secure engine based on the server configuration.
//
// Parameters:
// - cfg: *ServerCfg Configuration data for the HTTP server.
//
// Returns:
// - *Server: A new HTTP server.
// - error: An error if any.
func NewServer(cfg *ServerCfg) *Server {
	server := &Server{
		router: api.MakeRouter(),
		standard: &Engine{
			PORT:   cfg.HTTP,
			DOMAIN: cfg.DOMAIN,
			SUBS:   cfg.SUBS,
			server: newServerConfig(nil),
		},
	}

	if cfg.HTTPS == 0 {
		return server
	}

	server.secure = &Engine{
		PORT:   cfg.HTTPS,
		DOMAIN: cfg.DOMAIN,
		SUBS:   cfg.SUBS,
		server: newServerConfig(certs.TLSConfigFor(cfg.DOMAIN, cfg.SUBS...)),
	}

	http2.ConfigureServer(server.secure.server, &http2.Server{
		IdleTimeout: config.DEFAULT_TIMEOUT * time.Second,
	})

	return server
}

// Register registers API handlers with the server.
// This method binds API interfaces to the server for handling different routes and requests.
//
// Parameters:
// - api: api.APInterface An interface containing API handlers.
func (s *Server) Register(api api.APInterface) {
	// TODO: Implement registration logic.
}

// Start starts the HTTP server, allowing it to accept incoming connections.
// It checks for any running instances of the server and starts the standard and secure engines.
//
// Returns:
// - error: An error if any.
func (s *Server) Start() error {
	multi := errors.NewMultiError()

	if s.standard.running {
		multi.Add(errors.New("standard server active"))
	}

	if s.secure != nil && s.secure.running {
		multi.Add(errors.New("secure server active"))
	}

	if multi.Count() > 0 {
		return multi
	}

	multi.Add(s.standard.Start())
	if s.secure != nil {
		multi.Add(s.secure.Start())
	}

	return multi.IsError()
}

// Stop stops the HTTP server.
// This method halts the server's operations and stops accepting incoming connections.
//
// Returns:
// - error: An error if any.
func (s *Server) Stop() error {
	multi := errors.NewMultiError()

	if !s.standard.running {
		multi.Add(errors.New("standard server is not active"))
	}

	if s.secure != nil && !s.secure.running {
		multi.Add(errors.New("secure server is not active"))
	}

	if multi.Count() > 0 {
		return multi
	}

	multi.Add(s.standard.Stop())
	if s.secure != nil {
		multi.Add(s.secure.Stop())
	}

	return multi.IsError()
}

// Start starts the HTTP engine, allowing it to accept incoming connections.
// This method initiates the listening process on the specified port and handles incoming requests.
//
// Returns:
// - error: An error if any.
func (e *Engine) Start() error {
	if e.running {
		return errors.New("server already started")
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(e.PORT))
	if err != nil {
		return err
	}

	e.listener = listener
	if e.server.TLSConfig != nil {
		e.listener = tls.NewListener(e.listener, e.server.TLSConfig)
	}

	e.running = true

	logger.Info(fmt.Sprintf("server start on %v:%v with pid: %v", e.DOMAIN, e.PORT, os.Getpid()))

	go e.server.Serve(e.listener)

	return nil
}

// Stop stops the HTTP engine.
// It terminates the listening process and stops the server from accepting new connections.
//
// Returns:
// - error: An error if any.
func (e *Engine) Stop() error {
	if !e.running {
		return errors.New("server is not active")
	}

	err := e.listener.Close()
	if err != nil {
		e.running = false
	}

	logger.Info(fmt.Sprintf("server stop on %v:%v with pid: %v", e.DOMAIN, e.PORT, os.Getpid()))

	return err
}
