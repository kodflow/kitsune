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
	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/core/server/transport/promise"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

// Service struct represents a remote service to interact with.
type Service struct {
	Name         string // The name of the service, usually in URI form.
	Address      string // The address of the server
	Protocol     string // The network protocol to use (e.g., TCP, UDP)
	ID           string // A unique identifier for this connection
	Connected    bool   // True if a connection has been established, false otherwise
	tryReconnect bool
	network      net.Conn // The underlying network connection
}

// Create initializes a Service instance.
// address: The address of the remote service.
// port: The port number of the remote service.
// returns a pointer to a Service instance and an error if any.
func Create(address, port string) (*Service, error) {
	v4, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	s := &Service{
		Name:     address + ":" + port,
		Address:  address,
		Protocol: "tcp",
		ID:       v4.String(),
	}

	if err := s.Connect(); err != nil {
		return nil, err
	}

	return s, nil
}

// Connect establishes a connection to the server.
// Returns an error if the connection fails or if already connected.
func (s *Service) Connect() error {
	if s.Connected {
		return errors.New("already connected")
	}

	var err error
	s.network, err = net.DialTimeout(s.Protocol, s.Name, time.Second*config.DEFAULT_TIMEOUT)
	if err != nil {
		return fmt.Errorf("can't establish connection: %w", err)
	}

	s.Connected = true

	go s.handleServerResponses()

	return nil
}

// Disconnect closes the connection.
// Returns an error if the disconnection fails.
func (s *Service) Disconnect() error {
	s.Connected = false

	if err := s.network.Close(); err != nil {
		s.Connected = true
		return err
	}

	return nil
}

// Write sends data to the service.
// data: Buffer containing the data to be sent.
// returns the number of bytes written and an error if any.
func (s *Service) Write(data bytes.Buffer) (int, error) {
	if s.Connected {
		i, err := s.network.Write(data.Bytes())
		if err != nil {
			s.Connected = false
			s.reconnect()
			return 0, errors.New("lost connection")
		}

		return i, err
	}

	return 0, errors.New("not connected")
}

// reconnect tries to re-establish the connection every 5 seconds.
func (s *Service) reconnect() {
	if s.tryReconnect {
		return
	}

	s.tryReconnect = true
	defer func() { s.tryReconnect = false }()
	timeout := 0

	for {
		if s.Connected {
			return
		}

		if err := s.Connect(); err == nil {
			fmt.Println("Reconnected")
			return
		}

		time.Sleep(time.Second)
		timeout++
		if timeout > 10 {
			return
		}
	}
}

// handleServerResponses listens for responses from the server and processes them.
func (s *Service) handleServerResponses() {
	reader := bufio.NewReader(s.network)
	for {
		if !s.Connected {
			break
		}

		var length uint32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			break
		}

		data := make([]byte, length)
		_, err := io.ReadFull(reader, data)
		if err != nil {
			s.Disconnect()
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
			logger.Info(data)
			continue
		}

		if res.Pid != "" {
			p, err := promise.Find(res.Pid)
			if err != nil {
				fmt.Println(err)
				continue
			}

			p.Resolve(res)
		}
	}
}

func (s *Service) MakeExchange(answer ...bool) *Exchange {
	if len(answer) == 0 {
		return NewExchange(s.Name, true)
	}

	return NewExchange(s.Name, answer[0])
}

/*
func (s *Service) MakeCommand(answer ...bool) *cqrs.Message {
	if len(answer) == 0 {
		return cqrs.NewMessage(s.Name, true, cqrs.CMD)
	}

	return cqrs.NewMessage(s.Name, answer[0], cqrs.CMD)
}

/*
// MakeRequestWithResponse initializes a Query object with the expectation of a response.
// returns a Query instance configured to expect a response.
func (s *Service) MakeRequestWithResponse() *cqrs.Message {
	return cqrs.NewMessage(s.Name, true)
}

// MakeRequestOnly initializes a Query object without the expectation of a response.
// returns a Query instance configured to not expect a response.
func (s *Service) MakeRequestOnly() *cqrs.Message {
	return cqrs.NewMessage(s.Name, false)
}
*/
