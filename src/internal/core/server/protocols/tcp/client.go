package tcp

import (
	"fmt"
	"sync"

	"github.com/kodflow/kitsune/src/internal/core/server/transport"
)

// Client manages multiple service connections.
// This struct is responsible for managing connections to different services identified by their addresses.
type Client struct {
	services map[string]*Service
	mu       sync.Mutex // Mutex for thread-safe operations on services map
}

// NewClient creates a new instance of Client.
// It initializes the client with an empty map for managing services.
//
// Returns:
// - *Client: Newly created instance of Client.
func NewClient() *Client {
	return &Client{
		services: make(map[string]*Service),
	}
}

// Connect establishes a service connection for a given address.
// It creates an instance of Service for network interactions.
//
// Parameters:
// - address: string The TCP address to connect to.
//
// Returns:
// - error: Error, if any occurred during the connection setup.
func (c *Client) Connect(address string, nbInstances ...int) ([]*Service, error) {
	num := 1 // Default to one instance
	if len(nbInstances) > 0 && nbInstances[0] > 0 {
		num = nbInstances[0]
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if services, exists := c.services[address]; exists {
		return services, nil
	}

	services := make([]*Service, num)
	for i := 0; i < num; i++ {
		service := &Service{Address: address}
		if err := service.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect to service at %s: %v", address, err)
		}
		services[i] = service
	}

	c.services[address] = services
	return services, nil
}

// Send sends a request using the specified service and waits for a response.
// It uses the Service.Send method for the actual network operation.
//
// Parameters:
// - address: string The address of the service to send the request to.
// - exchange: *transport.Exchange The exchange object containing request and response.
//
// Returns:
// - *transport.Exchange: Exchange object containing the response.
// - error: Error, if any occurred during the send operation.
func (c *Client) Send(address string, exchange *transport.Exchange) (*transport.Exchange, error) {
	c.mu.Lock()
	service, exists := c.services[address]
	c.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("no service found for address %s", address)
	}

	return service.Send(exchange), nil
}

// Close closes all service connections managed by the client.
// It iterates over all services and closes each connection.
//
// Returns:
// - error: Error, if any occurred during the closure of services.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for address, service := range c.services {
		if err := service.Close(); err != nil {
			fmt.Printf("Error closing service at address %s: %v\n", address, err)
		}
	}

	// Clear the services map after closure
	c.services = make(map[string]*Service)

	return nil
}
