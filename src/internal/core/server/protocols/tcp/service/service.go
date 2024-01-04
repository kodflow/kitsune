// Package service provides the functionality to interact with a service over a network.
package service

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/core/server/protocols/tcp/promise"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

// Service struct represents a remote service to interact with.
type Service struct {
	Address      string   // The address of the service, usually in URI form.
	Protocol     string   // The network protocol to use (e.g., TCP, UDP).
	ID           string   // A unique identifier for this connection.
	Connected    bool     // True if a connection has been established, false otherwise.
	tryReconnect bool     // Flag indicating whether a reconnection attempt is in progress.
	Network      net.Conn // The underlying network connection.
}

// Create initializes a Service instance.
//
// Parameters:
// - address: The address of the remote service.
// - port: The port number of the remote service.
//
// Returns:
// - service: A pointer to a Service instance.
// - err: An error if any.
func Create(address string) (*Service, error) {
	v4, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	s := &Service{
		Address:  address,
		Protocol: "tcp",
		ID:       v4.String(),
	}

	return s.Connect()
}

// Connect establishes a connection to the server.
//
// Returns:
// - err: An error if the connection fails or if already connected.
func (s *Service) Connect() (*Service, error) {
	if s.Connected {
		return nil, errors.New("already connected")
	}

	var err error
	s.Network, err = net.DialTimeout(s.Protocol, s.Address, time.Second*config.DEFAULT_TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("can't establish connection: %w", err)
	}

	s.Connected = true

	go s.handleServerResponses()

	return s, nil
}

// Disconnect closes the connection.
//
// Returns:
// - err: An error if the disconnection fails.
func (s *Service) Disconnect() error {
	s.Connected = false

	if err := s.Network.Close(); err != nil {
		s.Connected = false
		return err
	}

	return nil
}

// Write sends data to the service.
//
// Parameters:
// - data: Buffer containing the data to be sent.
//
// Returns:
// - n: The number of bytes written.
// - err: An error if any.
func (s *Service) Write(data *bytes.Buffer) (int, error) {
	if s.Connected {
		n, err := s.Network.Write(data.Bytes())
		if err != nil {
			s.Connected = false
			s.reconnect()
			return 0, errors.New("lost connection")
		}

		return n, err
	}

	return 0, errors.New("not connected")
}

// reconnect tries to re-establish the connection every X seconds.
func (s *Service) reconnect() {
	// Check if a reconnection attempt is already in progress.
	if s.tryReconnect {
		return
	}

	// Set the tryReconnect flag to indicate a reconnection attempt is starting.
	s.tryReconnect = true

	// Use a deferred function to reset the tryReconnect flag when the function exits.
	defer func() { s.tryReconnect = false }()

	// Initialize a timeout counter.
	timeout := 0

	// Loop until the connection is re-established or until a timeout is reached.
	for {
		// If already connected, exit the loop.
		if s.Connected {
			return
		}

		// Try to establish a connection.
		if _, err := s.Connect(); err == nil {
			// Reconnection successful.
			fmt.Println("Reconnected")
			return
		}

		// Sleep for one second before attempting to reconnect again.
		time.Sleep(time.Second)
		timeout++
	}
}

// handleServerResponses listens to server responses and processes them.
func (s *Service) handleServerResponses() {
	// Create a reader to read data from the network connection.
	reader := bufio.NewReader(s.Network)

	// Continuously listen for server responses while the service is connected.
	for s.isConnected() {
		// Read the length of the incoming message.
		length, err := s.readMessageLength(reader)
		if logger.Error(err) {
			break // Exit the loop if there's an error reading the message length.
		}

		// Read the data with the specified length.
		data, err := s.readData(reader, length)
		if logger.Error(err) {
			// Handle any read errors and exit the loop.
			handleReadError(s, err)
			break
		}

		// Unmarshal the received data into a response object.
		res, err := unmarshalResponse(data)
		if logger.Error(err) {
			// Log an error if unmarshaling fails and continue to the next iteration.
			logger.Info(data)
			continue
		}

		// Process the received response.
		s.processResponse(res)
	}
}

// isConnected checks if the service is connected.
//
// Returns:
// - connected: True if the service is connected, false otherwise.
func (s *Service) isConnected() bool {
	return s.Connected
}

// readMessageLength reads the length of a message from the reader.
//
// Parameters:
// - reader: The bufio.Reader to read from.
//
// Returns:
// - length: The length of the message.
// - err: An error if any.
func (s *Service) readMessageLength(reader *bufio.Reader) (uint32, error) {
	var length uint32
	err := binary.Read(reader, binary.LittleEndian, &length)
	return length, err
}

// readData reads data of a specified length from the reader.
//
// Parameters:
// - reader: The bufio.Reader to read from.
// - length: The length of data to read.
//
// Returns:
// - data: The read data.
// - err: An error if any.
func (s *Service) readData(reader *bufio.Reader, length uint32) ([]byte, error) {
	data := make([]byte, length)
	_, err := io.ReadFull(reader, data)
	return data, err
}

// handleReadError handles read errors, disconnects the service if needed.
//
// Parameters:
// - s: The Service instance.
// - err: The error to handle.
func handleReadError(s *Service, err error) {
	s.Disconnect()
	if err == io.EOF {
		logger.Info("connection closed by the server.")
	} else {
		logger.Error(err)
	}
}

// unmarshalResponse unmarshals a response from binary data.
//
// Parameters:
// - data: The binary data to unmarshal.
//
// Returns:
// - res: The unmarshaled response.
// - err: An error if any.
func unmarshalResponse(data []byte) (*generated.Response, error) {
	res := &generated.Response{}
	err := proto.Unmarshal(data, res)
	return res, err
}

// processResponse processes a response, resolves associated promises.
//
// Parameters:
// - res: The response to process.
func (s *Service) processResponse(res *generated.Response) {
	if res.Pid != "" {
		p, err := promise.Find(res.Pid)
		if logger.Error(err) {
			return
		}

		p.Resolve(res)
	}
}

// MakeExchange creates a new Exchange instance for this service.
//
// Parameters:
// - answer: Optional boolean argument to specify if the exchange should be answered.
//
// Returns:
// - exchange: A pointer to a new Exchange instance.
func (s *Service) MakeExchange(answer ...bool) *Exchange {
	if len(answer) == 0 {
		return NewExchange(s.Address, true)
	}

	return NewExchange(s.Address, answer[0])
}
