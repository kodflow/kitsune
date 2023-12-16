package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kodmain/kitsune/src/internal/core/certs"
	"github.com/kodmain/kitsune/src/internal/core/server/api"
	"github.com/kodmain/kitsune/src/internal/core/server/handler"
	"github.com/kodmain/kitsune/src/internal/kernel/errors"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"golang.org/x/net/http2"
)

type Server struct {
	standard *Engine
	secure   *Engine
	router   *api.Router
}

type ServerCfg struct {
	DOMAIN string
	SUBS   []string
	HTTP   int
	HTTPS  int
}

type Engine struct {
	PORT     int
	DOMAIN   string
	SUBS     []string
	listener net.Listener
	server   *http.Server
	running  bool
}

func NewServer(cfg *ServerCfg) *Server {
	server := &Server{
		router: api.MakeRouter(),
		standard: &Engine{
			PORT:   cfg.HTTP,
			DOMAIN: cfg.DOMAIN,
			SUBS:   cfg.SUBS,
			server: &http.Server{
				Handler:           http.HandlerFunc(handler.HTTPHandler),
				ReadTimeout:       5 * time.Second,
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       120 * time.Second,
				ReadHeaderTimeout: 2 * time.Second,
				MaxHeaderBytes:    1 << 20,
			},
		},
	}

	if cfg.HTTPS == 0 {
		return server
	}

	server.secure = &Engine{
		PORT:   cfg.HTTPS,
		DOMAIN: cfg.DOMAIN,
		SUBS:   cfg.SUBS,
		server: &http.Server{
			Handler:           http.HandlerFunc(handler.HTTPHandler),
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       120 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			MaxHeaderBytes:    1 << 20,
			TLSConfig:         certs.TLSConfigFor(cfg.DOMAIN, cfg.SUBS...),
		},
	}

	http2.ConfigureServer(server.secure.server, &http2.Server{
		IdleTimeout: 120 * time.Second,
	})

	return server
}

// Register is a method for registering API handlers with the server.
func (s *Server) Register(api api.APInterface) {
	// TODO
}

// Start starts the TCP server, allowing it to accept incoming connections.
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

// Stop stops the TCP server.
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

	logger.Info("server stop on " + s.standard.DOMAIN)

	return multi.IsError()
}

// Start starts the TCP server, allowing it to accept incoming connections.
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

	logger.Info("server start on " + strconv.Itoa(e.PORT) + " with pid:" + strconv.Itoa(os.Getpid()))

	go e.server.Serve(e.listener)

	return nil
}

// Stop stops the TCP server.
func (s *Engine) Stop() error {
	if !s.running {
		return errors.New("server is not active")
	}

	err := s.listener.Close()
	if err != nil {
		s.running = false
	}

	logger.Info("server stop on ")
	return err
}
