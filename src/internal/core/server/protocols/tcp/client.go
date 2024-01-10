package tcp

import (
	"fmt"
	"sync"

	"github.com/kodflow/kitsune/src/config"
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
// - nbInstances: ...int Optional parameter to specify the number of instances to create (default is 1).
//
// Returns:
// - *Service: Instance of the Service for the specified address.
// - error: Error, if any occurred during the connection setup.
func (c *Client) Connect(address string, nbInstances ...int) (*Service, error) {
	num := config.DEFAULT_CLIENT_SERVICE_MAX_CONNS // Default to one instance
	if len(nbInstances) > 0 && nbInstances[0] > 0 {
		num = nbInstances[0]
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if service, exists := c.services[address]; exists {
		return service, nil
	}

	c.services[address] = NewService(address, num)

	return c.services[address], nil
}

// Close closes all service connections managed by the client.
// It iterates over all services and closes each connection.
//
// Returns:
// - error: Error, if any occurred during the closure of services.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, service := range c.services {
		if err := service.Close(); err != nil {
			fmt.Printf("Error closing service at address %s: %v\n", service.address, err)
		}
	}

	// Clear the services map after closure
	c.services = make(map[string]*Service)

	return nil
}
