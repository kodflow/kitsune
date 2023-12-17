// Package tcp provides functionalities for a TCP client.
package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/kodflow/kitsune/src/internal/core/server/transport/promise"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/service"
	"google.golang.org/protobuf/proto"
)

// Client represents a TCP client for interacting with remote services.
type Client struct {
	mu       sync.Mutex                  // Mutex to protect the services map
	services map[string]*service.Service // A map of services by name
}

// NewClient creates a new TCP client.
func NewClient() *Client {
	return &Client{
		services: make(map[string]*service.Service),
	}
}

// Connect establishes a connection to a remote service.
//
// Parameters:
// - address: string The address of the remote service.
// - port: string The port of the remote service.
//
// Returns:
// - s: *service.Service The connected service.
// - error: An error if any.
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

// Disconnect closes connections to one or more services.
//
// Parameters:
// - services: ...string A variadic list of service names to disconnect from.
//
// Returns:
// - error: An error if any.
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

// Send sends requests to remote services and handles the responses.
//
// Parameters:
// - callback: func(...*generated.Response) A function to handle responses.
// - queries: ...*service.Exchange An array of service exchanges containing requests and their associated services.
//
// Returns:
// - error: An error if any.
func (c *Client) Send(callback func(...*generated.Response), queries ...*service.Exchange) error {
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

// validateInputs checks if the inputs to the Send function are valid.
//
// Parameters:
// - queries: []*service.Exchange An array of service exchanges containing requests and their associated services.
//
// Returns:
// - error: An error if any.
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

// buildDispatch organizes queries by service name for dispatch.
//
// Parameters:
// - queries: []*service.Exchange An array of service exchanges containing requests and their associated services.
//
// Returns:
// - dispatch: map[string][]*service.Exchange A map of service names to their respective queries.
func (c *Client) buildDispatch(queries []*service.Exchange) map[string][]*service.Exchange {
	dispatch := make(map[string][]*service.Exchange)

	for _, query := range queries {
		if _, ok := c.services[query.Service]; ok {
			dispatch[query.Service] = append(dispatch[query.Service], query)
		}
	}

	return dispatch
}

// processQueries prepares queries for sending and sets up callback handling.
//
// Parameters:
// - dispatch: map[string][]*service.Exchange A map of service names to their respective queries.
// - callback: func(...*generated.Response) A function to handle responses.
//
// Returns:
// - p: *promise.Promise A promise object for handling responses.
// - error: An error if any.
func (c *Client) processQueries(dispatch map[string][]*service.Exchange, callback func(...*generated.Response)) (*promise.Promise, error) {
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

// sendToServices sends queries to remote services.
//
// Parameters:
// - dispatch: map[string][]*service.Exchange A map of service names to their respective queries.
//
// Returns:
// - error: An error if any.
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
