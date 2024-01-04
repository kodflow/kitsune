package http

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/core/certs"
	"github.com/kodflow/kitsune/src/internal/core/server/router"
	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/kernel/errors"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"golang.org/x/net/http2"
)

// Server represents an HTTP server that handles both standard and secure connections.
// It encapsulates the functionality of two engines, one for handling standard HTTP
// connections and the other for secure HTTPS connections, along with a router for API routing.
type Server struct {
	standard *Engine        // The engine for handling standard HTTP connections.
	secure   *Engine        // The engine for handling secure HTTPS connections.
	router   *router.Router // The router for managing API endpoints.
}

// ServerCfg holds configuration data for the HTTP server.
// This structure includes details such as the server's domain, subdomains, and port numbers
// for both HTTP and HTTPS connections.
type ServerCfg struct {
	DOMAIN string   // The domain of the server.
	SUBS   []string // The subdomains of the server.
	HTTP   string   // The port number for HTTP connections.
	HTTPS  string   // The port number for HTTPS connections.
}

// Engine represents an HTTP engine.
// It is responsible for managing network listeners and the HTTP server for either standard
// or secure connections, keeping track of its running status and configuration details
// like port, domain, and subdomains.
type Engine struct {
	PORT     string       // The port number for the engine.
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
		router: router.MakeRouter(),
		standard: &Engine{
			PORT:   cfg.HTTP,
			DOMAIN: cfg.DOMAIN,
			SUBS:   cfg.SUBS,
			server: newServerConfig(nil),
		},
	}

	server.standard.server.Handler = http.HandlerFunc(server.HTTPHandler)

	if cfg.HTTPS == "" {
		return server
	}

	server.secure = &Engine{
		PORT:   cfg.HTTPS,
		DOMAIN: cfg.DOMAIN,
		SUBS:   cfg.SUBS,
		server: newServerConfig(certs.TLSConfigFor(cfg.DOMAIN, cfg.SUBS...)),
	}

	server.secure.server.Handler = http.HandlerFunc(server.HTTPHandler)

	http2.ConfigureServer(server.secure.server, &http2.Server{
		IdleTimeout: config.DEFAULT_TIMEOUT * time.Second,
	})

	return server
}

func (s *Server) Register(api *router.EndPoint) {
	s.router.Register(api)
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

	listener, err := net.Listen("tcp", ":"+e.PORT)
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

// HTTPHandler handles HTTP requests and sends back HTTP responses.
// It processes incoming HTTP requests, creates a corresponding transport request,
// and uses the router to generate a response. The handler deals with various HTTP methods,
// reads request bodies if necessary, and writes back responses including headers and status codes.
//
// Parameters:
// - w: http.ResponseWriter Response writer to send back the HTTP response.
// - r: *http.Request The incoming HTTP request to be processed.
func (s *Server) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info(fmt.Sprintf("request received on %v:%v with pid: %v", r.Host, r.URL.String(), os.Getpid()))
	// Initialize a new transport request and response
	req, res := transport.New()
	req.Method = r.Method
	req.Endpoint = r.URL.String()

	// Read the request body for specific HTTP methods
	if r.Method == "POST" || r.Method == "PATCH" || r.Method == "PUT" {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			// Handle errors in reading the request body
			http.Error(w, "Erreur lors de la lecture de la requête", http.StatusBadRequest)
			return
		}

		req.Body = body
	}

	// Process the request using the router
	if err := s.router.Resolve(req, res); err != nil {
		// Handle errors in processing the request
		http.Error(w, "Erreur de traitement", http.StatusInternalServerError)
		return
	}

	// Write the response back to the client
	w.Header().Set("request-id", req.Id)
	w.WriteHeader(int(res.Status))
	w.Write(res.Body)
}
