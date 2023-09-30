// Package socket provides functionalities for both a TCP client and a TCP server.
// It enables the creation, management, and communication between clients and the server over TCP.
// Messages sent between the client and server are serialized using protobuf.

// socket provides functionalities for a TCP client.
package socket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/kodmain/kitsune/src/internal/core/server/promise"
	"github.com/kodmain/kitsune/src/internal/core/server/service"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"google.golang.org/protobuf/proto"
)

/*
func init() {
	var address, port, protocol string = "", "", ""

	client := NewClient()
	service1, _ := client.Connect(address, port, protocol)
	service2, _ := client.Connect(address, port, protocol)
	service3, _ := client.Connect(address, port, protocol)

	req1 := service1.MakeRequestWithResponse()
	req2 := service2.MakeRequestWithResponse()
	req3 := service2.MakeRequestWithResponse()
	req4 := service3.MakeRequestOnly()

	client.Send(func(responses ...*transport.Response) {
		fmt.Println(responses)
	}, req1, req2, req3, req4)
}
*/

// Client represents a TCP client with functionalities such as sending requests and waiting for responses.
type Client struct {
	mu       sync.Mutex
	services map[string]*service.Service
}

// NewClient initializes and returns a new Client instance.
// address is the TCP address for the client.
func NewClient() *Client {
	c := &Client{
		services: map[string]*service.Service{},
	}

	return c
}

func (c *Client) Connect(address, port, protocol string) (*service.Service, error) {
	s, err := service.Create(address, port, protocol)
	if err != nil {
		return s, err
	}

	c.mu.Lock()
	c.services[s.Name] = s
	c.mu.Unlock()

	return s, nil
}

// Disconnect terminates the active connection if it exists.
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

// Send transmits a request to the server and returns a promise for the response.
// req is the request to be sent.
func (c *Client) Send(callback func(...*transport.Response), queries ...*service.Query) error {
	if len(queries) == 0 {
		return fmt.Errorf("no request")
	}

	dispatch := map[string][]*service.Query{}
	buffers := map[string]bytes.Buffer{}

	c.mu.Lock()
	services := c.services
	c.mu.Unlock()

	for _, query := range queries {
		if _, ok := services[query.Service]; ok {
			dispatch[query.Service] = append(dispatch[query.Service], query)
		}
	}

	p, err := promise.Create(callback)
	if err != nil {
		return err
	}

	for service, queries := range dispatch {
		var buffer bytes.Buffer

		for _, query := range queries {
			if query.Answer {
				p.Add(query.Req)
			}

			data, err := proto.Marshal(query.Req)
			if err != nil {
				return err
			}

			// Écrire la longueur du message
			if err := binary.Write(&buffer, binary.LittleEndian, uint32(len(data))); err != nil {
				return err
			}

			// Écrire les données du message
			if _, err := buffer.Write(data); err != nil {
				return err
			}

			buffers[service] = buffer
		}
	}

	for service, buffer := range buffers {
		services[service].Write(buffer)
	}

	return nil
}

/*
// SendSync envoie une requête de manière synchrone et attend la réponse avant de retourner.
func (c *Client) SendSync(req *transport.Request) (*transport.Response, error) {
	p, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if !req.Answer {
		return nil, errors.New("you send a non answer query")
	}

	// Attendre la réponse de manière synchrone
	res := p.Wait()
	if res == nil {
		return nil, errors.New("no response received")
	}

	return res, nil
}
*/
