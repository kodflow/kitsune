// server.go

package http

import (
	"context"
	"net/http"
	"runtime"

	"golang.org/x/net/http2"
)

// Server represents an HTTP server.
type Server struct {
	Address string             // Address on which the server will listen.
	ctx     context.Context    // Context to signal server termination.
	cancel  context.CancelFunc // Function to cancel the server context.
	server  *http.Server       // Internal HTTP server.
	tasks   chan func()        // Task channel for workers.
}

// NewV1 creates a new HTTP/1.1 server.
func NewV1(address string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{
		Address: address,
		ctx:     ctx,
		cancel:  cancel,
		server:  &http.Server{Addr: address},         // Initialisation du serveur HTTP ici.
		tasks:   make(chan func(), runtime.NumCPU()), // Initialize the task channel.
	}

	// Start workers.
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.worker()
	}

	return s
}

// NewV2 creates a new HTTP/2 server.
func NewV2(address string) *Server {
	s := NewV1(address)
	http2.ConfigureServer(s.server, &http2.Server{}) // Maintenant, s.server n'est plus nil.
	return s
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	// Implement the logic to start the server.
	return nil
}

// Stop stops the HTTP server.
func (s *Server) Stop() error {
	// Notify all workers to stop.
	close(s.tasks)
	s.cancel()
	// Implement the logic to stop the server.
	return nil
}

// worker is a goroutine that processes tasks.
func (s *Server) worker() {
	for {
		select {
		case task, ok := <-s.tasks:
			// If the tasks channel is closed, exit the worker.
			if !ok {
				return
			}
			// Run the task.
			task()
		// If context is done, exit the worker.
		case <-s.ctx.Done():
			return
		}
	}
}

/*
func (s *Server) worker() {
	for request := range requests {
		go s.handler()
	}
}

// StartServer initializes and starts the HTTP/2 server
func (s *Server) Start() {

	for i := 0; i < runtime.NumCPU()*config.DEFAULT_IO_BOUND; i++ {
		go s.worker()
	}

	http2.ConfigureServer(server, &http2.Server{})

	log.Printf("Starting HTTP/2 server on %s", address)
	if err := server.ListenAndServeTLS("path/to/cert.pem", "path/to/key.pem"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
*/
