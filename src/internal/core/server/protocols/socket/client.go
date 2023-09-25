// Package socket provides functionalities for both a TCP client and a TCP server.
// It enables the creation, management, and communication between clients and the server over TCP.
// Messages sent between the client and server are serialized using protobuf.

// socket provides functionalities for a TCP client.
package socket

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

// Client represents a TCP client with functionalities such as sending requests and waiting for responses.
type Client struct {
	Address   string                              // Address is the TCP address of the client.
	conn      net.Conn                            // conn is the active connection instance.
	responses map[string]chan *transport.Response // responses store channels for awaiting responses based on request ID.
	mu        sync.Mutex                          // mu is a mutex for handling concurrent access to the responses map.
}

// NewClient initializes and returns a new Client instance.
// address is the TCP address for the client.
func NewClient(address string) *Client {
	c := &Client{
		Address:   address,
		responses: make(map[string]chan *transport.Response),
	}

	return c
}

func (c *Client) Connect() error {
	if c.conn != nil {
		return fmt.Errorf("already connected")
	}

	var err error
	for i := 0; i < config.DEFAULT_RETRY_MAX; i++ {
		c.conn, err = net.DialTimeout("tcp", c.Address, time.Second*config.DEFAULT_TIMEOUT)

		if err == nil {
			go c.handleServerResponses()
			return nil
		}

		time.Sleep(config.DEFAULT_RETRY_INTERVAL)
	}

	return fmt.Errorf("failed to connect after %d attempts", config.DEFAULT_RETRY_MAX)
}

// Disconnect terminates the active connection if it exists.
func (c *Client) Disconnect() error {
	if c.conn == nil {
		return errors.New("connection already closed")
	}

	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.conn = nil

	return nil
}

// Send transmits a request to the server and returns a promise for the response.
// req is the request to be sent.
func (c *Client) Send(req *transport.Request) (*promise, error) {
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	p := Promise()
	if req.Answer {
		p.Init()
		c.mu.Lock()
		if _, exists := c.responses[req.Id]; exists {
			c.mu.Unlock()
			return nil, fmt.Errorf("request ID %s is already in use", req.Id)
		}
		c.responses[req.Id] = p.wait
		c.mu.Unlock()
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil, errors.New("connection is closed")
	}

	if err := binary.Write(c.conn, binary.LittleEndian, uint32(len(data))); err != nil {
		return nil, err
	}

	_, err = c.conn.Write(data)
	if err != nil {
		return nil, err
	}

	return p, nil
}

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

// handleServerResponses continuously reads responses from the server and forwards them to the appropriate channels.
func (c *Client) handleServerResponses() {
	reader := bufio.NewReader(c.conn)

	for {
		if c.conn == nil {
			break
		}

		var length uint32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			break
		}

		data := make([]byte, length)
		_, err := io.ReadFull(reader, data)
		if err != nil {
			if err == io.EOF {
				logger.Info("connection closed by the server.")
			} else {
				logger.Error(err)
			}
			break
		}

		res := &transport.Response{}
		err = proto.Unmarshal(data, res)
		if logger.Error(err) {
			continue
		}

		c.mu.Lock()
		if ch, exists := c.responses[res.Id]; exists {
			ch <- res
			close(ch)
			delete(c.responses, res.Id)
		}
		c.mu.Unlock()
	}
}
