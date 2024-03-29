package tcp

import (
	"sync"

	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
)

type Service struct {
	mutex       sync.Mutex    // Mutex for thread-safe access.
	connections []*Connection // Active TCP connections.
	address     string        // TCP server address.
	current     int           // Current number of connections.

	recover  chan *generated.Response
	promises map[string]*transport.Exchange
}

// NewService creates a new service instance.
// Initializes the service and starts connection cleanup routine.
//
// Parameters:
// - address: string The TCP address of the server.
// - maxConns: int Maximum number of connections.
//
// Returns:
// - *Service: New service instance.
func NewService(address string, maxConns int) *Service {
	service := &Service{
		address:  address,
		recover:  make(chan *generated.Response),
		promises: make(map[string]*transport.Exchange),
	}

	for i := 0; i < maxConns; i++ {
		service.connections = append(service.connections, newConnection(address, service.recover))
	}

	go service.aggregate()

	return service
}

func (s *Service) aggregate() {
	for p := range s.recover {
		s.mutex.Lock()
		s.promises[p.Id].Response(p)
		delete(s.promises, p.Id)
		s.mutex.Unlock()
	}
}

// Send sends a request and waits for a response.
// Uses an available connection from the pool or scales up if needed.
//
// Parameters:
// - exchange: *transport.Exchange Exchange object with request and response.
//
// Returns:
// - *transport.Exchange: Updated exchange object with response.
func (s *Service) Send(exchange *transport.Exchange) *transport.Exchange {
	return s.process(exchange, s.current%len(s.connections))
}

// process the request using a specific connection.
// It uses the indexed writer and reader for sending the request and receiving the response.
//
// Parameters:
// - exchange: *transport.Exchange The exchange object containing the request and response.
// - index: int The index of the connection to use for this request.
//
// Returns:
// - *transport.Exchange: The exchange object with the updated response.
func (s *Service) process(exchange *transport.Exchange, index int) *transport.Exchange {
	req := exchange.Request()

	s.mutex.Lock()
	conn := s.connections[index]
	s.promises[req.Id] = exchange
	s.current = (s.current + 1) % len(s.connections)
	s.mutex.Unlock()

	conn.o <- req

	return exchange
}

// Close closes all TCP connections of the service.
// It stops the cleanup ticker and closes each active connection in the service.
//
// Returns:
// - error: An error, if any occurred during the closure of connections.
func (s *Service) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var err error

	for i, conn := range s.connections {
		if conn != nil {
			conn.mutex.Lock()
			if conn.net != nil {
				conn.close = true
				if closeErr := conn.net.Close(); closeErr != nil {
					err = closeErr // Set the error if closing a connection fails
				}
				conn.net = nil
			}
			s.connections[i] = nil
			conn.mutex.Unlock()
		}
	}

	return err
}
