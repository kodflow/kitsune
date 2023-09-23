// Package socket provides functionalities for both a TCP client and a TCP server.
// It enables the creation, management, and communication between clients and the server over TCP.
// Messages sent between the client and server are serialized using protobuf.

// socket provides functionalities for a TCP server.
package socket

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/core/server/router"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

// Server represents a TCP server with the capability to manage multiple clients.
// It contains a listener to accept incoming connections, and a concurrent map to keep track of connected clients.
type Server struct {
	Address  string             // Address specifies the TCP address for the server to listen on.
	ctx      context.Context    // ctx represents the server's context to manage its lifecycle.
	cancel   context.CancelFunc // cancel function to signal the termination of the server's operations.
	listener net.Listener       // listener is the actual TCP listener for the server.
}

// NewServer initializes a new server instance with the specified listening address.
// It returns an instance of the Server.
func NewServer(address string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		Address: address,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start initiates the server to begin listening for incoming connections.
// It spins up worker goroutines equivalent to the number of available CPUs and starts the listener.
// Returns an error if the server is already running or if there's an issue starting the listener.
func (s *Server) Start() error {
	if s.listener != nil {
		return errors.New("server already started")
	}

	var err error
	s.listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	connections := make(chan net.Conn, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU()*config.DEFAULT_IO_BOUND; i++ {
		go s.worker(connections)
	}

	go func() {
		defer close(connections)
		for {
			if s.listener == nil {
				break
			}

			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}

			connections <- conn
		}
	}()

	return nil
}

// Stop terminates the server's operations, closes the listener, and clears the client map.
// Returns an error if there's an issue closing the listener or if the server is not active.
func (s *Server) Stop() error {
	if s.listener == nil {
		return errors.New("server is not active")
	}

	err := s.listener.Close()
	s.listener = nil
	return err
}

// worker is a dedicated goroutine that serves incoming connections.
// It reads from the connections channel and spawns a new goroutine to handle each connection.
func (s *Server) worker(conns chan net.Conn) {
	for conn := range conns {
		go s.handleConnection(conn)
	}
}

// handleConnection manages the lifecycle of a single client connection.
// It reads incoming data, processes it, and sends back responses.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	data, err := s.handleRequest(reader)
	if err != nil {
		logger.Error(err)
		return
	}

	req := &transport.Request{}
	err = proto.Unmarshal(data, req)
	if err != nil {
		logger.Error(err)
		return
	}

	cid := string(req.Body)

	for {
		data, err := s.handleRequest(reader)
		if err != nil {
			logger.Error(fmt.Errorf("server lost connection with %s", cid))
			break
		}
		go s.handleResponse(conn, data)
	}
}

// handleRequest processes an incoming request from a client connection.
// It reads the data, decodes it, and returns the data as bytes.
func (s *Server) handleRequest(reader *bufio.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err := io.ReadFull(reader, data)
	return data, err
}

// handleResponse sends a response back to the client after processing the incoming request.
// It encodes the data and writes it to the client connection.
func (s *Server) handleResponse(conn net.Conn, b []byte) {
	req := transport.RequestFromBytes(b)
	res := router.Resolve(req)
	if req.Answer {
		data, err := transport.ResponseToBytes(res)
		logger.Error(err)
		err = binary.Write(conn, binary.LittleEndian, uint32(len(data)))
		logger.Error(err)
		_, err = conn.Write(data)
		logger.Error(err)
	}
}
