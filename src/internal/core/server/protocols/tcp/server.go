// Package tcp provides functionalities for both a TCP client and a TCP server.
// It enables the creation, management, and communication between clients and the server over TCP.
// Messages sent between the client and server are serialized using protobuf.

// tcp provides functionalities for a TCP server.
package tcp

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/kodmain/kitsune/src/internal/core/server/service/router"
)

// Server represents a TCP server with the capability to manage multiple clients.
type Server struct {
	Address  string             // Address specifies the TCP address for the server to listen on.
	ctx      context.Context    // ctx represents the server's context to manage its lifecycle.
	cancel   context.CancelFunc // cancel function to signal the termination of the server's operations.
	listener net.Listener       // listener is the actual TCP listener for the server.
}

// NewServer initializes a new server instance.
// address: The listening address for the server.
func NewServer(address string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		Address: address,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start initiates the server to begin listening for incoming connections.
func (s *Server) Start() error {
	if s.listener != nil {
		return errors.New("server already started")
	}

	var err error
	s.listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	go func() {
		for {
			if s.listener == nil {
				break
			}

			conn, err := s.listener.Accept()
			if err != nil {
				break
			}

			go s.handleConnection(conn)
		}
	}()

	return nil
}

// Stop terminates the server's operations.
func (s *Server) Stop() error {
	if s.listener == nil {
		return errors.New("server is not active")
	}

	err := s.listener.Close()
	s.listener = nil
	return err
}

// handleConnection manages a single client connection.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		data, err := s.handleRequest(reader)
		if err != nil {
			break
		}
		s.sendResponse(writer, data)
	}
}

// handleRequest processes an incoming request.
// reader: A buffered reader for the incoming data.
func (s *Server) handleRequest(reader *bufio.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err := io.ReadFull(reader, data)
	return data, err
}

// sendResponse sends a response back to the client.
// conn: The client connection instance.
// b: The byte array containing the request.
func (s *Server) sendResponse(writer *bufio.Writer, b []byte) {
	res := router.Handler(b)
	if len(res) > 0 {
		binary.Write(writer, binary.LittleEndian, uint32(len(res)))
		writer.Write(res)
		writer.Flush()
	}
}
