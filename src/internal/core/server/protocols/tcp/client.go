// Package tcp provides functionalities for a TCP server.
package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/kodmain/kitsune/src/internal/core/server/service"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/core/server/transport/promise"
	"google.golang.org/protobuf/proto"
)

// Client represents a TCP client.
// It manages multiple service connections and provides methods to send and receive data.
type Client struct {
	mu       sync.Mutex                  // Protects the services map
	services map[string]*service.Service // A map of services by name
}

// NewClient initializes a new Client and returns its pointer.
//
// Returns:
//
//	*Client: A pointer to the newly created Client
func NewClient() *Client {
	return &Client{
		services: make(map[string]*service.Service),
	}
}

// Connect initiates a new service connection.
// Parameters:
//
//	address: The IP address of the service
//	port: The port number of the service
//
// Returns:
//
//	*service.Service: A pointer to the connected service
//	error: An error object if an error occurred
func (c *Client) Connect(address, port string) (*service.Service, error) {
	s, err := service.Create(address, port)
	if err != nil {
		return s, err
	}

	c.mu.Lock()
	c.services[s.Name] = s
	c.mu.Unlock()

	return s, nil
}

// Disconnect terminates active connections to specified services or all services if none are specified.
//
// Parameters:
//
//	services: Names of services to disconnect (variadic)
//
// Returns:
//
//	error: An error object if an error occurred
func (c *Client) Disconnect(services ...string) error {
	if len(c.services) == 0 {
		return errors.New("no connection")
	}

	if len(services) == 0 {
		for service, mp := range c.services {
			if err := mp.Disconnect(); err != nil {
				continue
			}
			delete(c.services, service)
		}

		return nil
	}

	for _, service := range services {
		delete(c.services, service)
	}

	return nil
}

// Send sends queries to services and registers a callback for responses.
//
// Parameters:
//
//	callback: A function to be called when responses are received
//	queries: A slice of queries to send to services (variadic)
//
// Returns:
//
//	error: An error object if an error occurred
//
// Send sends queries to services and registers a callback for responses.
func (c *Client) Send(callback func(...*transport.Response), queries ...*service.Exchange) error {
	if err := c.validateInputs(queries); err != nil {
		return err
	}

	dispatch := c.buildDispatch(queries)
	p, err := c.processQueries(dispatch, callback)
	if err != nil {
		return err
	}

	if err := c.sendToServices(dispatch); err != nil {
		return err
	}

	if p.Length == 0 {
		p.Close()
	}

	return nil
}

// validateInputs validates the service connections and queries.
func (c *Client) validateInputs(queries []*service.Exchange) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.services) == 0 {
		return errors.New("no connection")
	}

	if len(queries) == 0 {
		return fmt.Errorf("no request")
	}

	return nil
}

// buildDispatch builds the dispatch map for service names and their queries.
func (c *Client) buildDispatch(queries []*service.Exchange) map[string][]*service.Exchange {
	dispatch := make(map[string][]*service.Exchange)

	for _, query := range queries {
		if _, ok := c.services[query.ServiceName()]; ok {
			dispatch[query.ServiceName()] = append(dispatch[query.ServiceName()], query)
		}
	}

	return dispatch
}

// processQueries processes the queries and marshals the request data.
func (c *Client) processQueries(dispatch map[string][]*service.Exchange, callback func(...*transport.Response)) (*promise.Promise, error) {
	p, err := promise.Create(callback)
	if err != nil {
		return nil, err
	}

	for _, queries := range dispatch {
		for _, query := range queries {
			if query.Answer {
				p.Add(query.Req)
			}
		}
	}

	return p, nil
}

// sendToServices sends the marshaled data to the services.
func (c *Client) sendToServices(dispatch map[string][]*service.Exchange) error {
	buffers := make(map[string]*bytes.Buffer)

	for service, queries := range dispatch {
		var buffer bytes.Buffer

		for _, query := range queries {
			data, err := proto.Marshal(query.Req)
			if err != nil {
				return err
			}

			if err := binary.Write(&buffer, binary.LittleEndian, uint32(len(data))); err != nil {
				return err
			}

			if _, err := buffer.Write(data); err != nil {
				return err
			}
		}

		buffers[service] = &buffer
	}

	for service, buffer := range buffers {
		if _, err := c.services[service].Write(buffer); err != nil {
			return err
		}
	}

	return nil
}
