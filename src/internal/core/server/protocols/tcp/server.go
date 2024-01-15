package tcp

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/kodflow/kitsune/src/internal/core/server/router"
	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
)

// Server represents a TCP server and contains information about the address it listens on
// and the underlying network listener.
type Server struct {
	Address   string       // Address to listen on
	listener  net.Listener // TCP Listener object
	router    *router.Router
	isRunning bool

	i chan []byte
	o chan []byte
}

// NewServer creates a new Server instance with the specified listening address.
func NewServer(address string) *Server {
	return &Server{
		Address: address,
		router:  router.MakeRouter(),
		// I/O channels
		i: make(chan []byte),
		o: make(chan []byte),
	}
}

// Register is a method for registering API handlers with the server.
//
// Parameters:
// - api: router.EndPoint - The EndPoint to register handlers from.
func (s *Server) Register(api router.EndPoint) {
	// TODO: Implement handler registration.
}

// Start starts the TCP server, allowing it to accept incoming connections.
//
// Returns:
// - error: An error if the server is already started or if there was an issue starting the server.
func (s *Server) Start() error {
	if s.listener != nil {
		return errors.New("server already started")
	}

	var err error
	s.listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	go s.acceptLoop()

	logger.Info("server start on " + s.Address + " with pid:" + strconv.Itoa(os.Getpid()))

	return nil
}

// Stop stops the TCP server.
//
// Returns:
// - error: An error if the server is not active or if there was an issue stopping the server.
func (s *Server) Stop() error {
	if s.listener == nil {
		return errors.New("server is not active")
	}

	err := s.listener.Close()
	s.listener = nil
	s.isRunning = false

	logger.Info("server stop on " + s.Address)
	return err
}

// accepLoop continuously accepts incoming connections.
// It listens for incoming client connections and handles them asynchronously by calling 'handleConnection'.
func (s *Server) acceptLoop() {
	s.isRunning = true
	for {
		if s.listener == nil {
			break
		}

		conn, err := s.listener.Accept() // Accept incoming connections.
		if !s.isRunning || err != nil {
			break
		}

		go s.handleConnection(conn) // Handle the connection asynchronously using 'handleConnection'.
	}
}

// handleConnection handles incoming client connections.
// It reads data from the 'conn', processes incoming requests, and sends responses back.
//
// Parameters:
// - conn: net.Conn - The client connection to handle.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	go s.request(reader)
	go s.response(writer)

	for data := range s.i {
		go s.TCPHandler(data)
	}
}

func (s *Server) request(reader *bufio.Reader) {
	for {
		var length uint32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			logger.Error(fmt.Errorf("failed to read request length: %w", err))
			continue
		}

		data := make([]byte, length)
		_, err := io.ReadFull(reader, data)

		if err != nil {
			logger.Error(fmt.Errorf("failed to read request: %w", err))
			continue
		}
		s.i <- data
	}
}

func (s *Server) response(writer *bufio.Writer) {
	for data := range s.o {
		if len(data) > 0 {
			binary.Write(writer, binary.LittleEndian, uint32(len(data)))
			writer.Write(data)
			writer.Flush()
		}
	}
}

// TCPHandler handles TCP requests by unmarshalling, processing, and marshalling responses.
// It is responsible for converting raw byte data into a structured request, processing it
// using a router, and then returning the structured response as byte data. It handles
// errors at each step by returning an empty response in case of failure.
//
// Parameters:
// - b: []byte Raw byte array representing a TCP request.
//
// Returns:
// - []byte: Processed response as a byte array. Returns an empty response in case of errors.
func (s *Server) TCPHandler(b []byte) {
	exchange := transport.New()
	exchange.RequestFromTCP(b)
	s.router.Resolve(exchange)
	s.o <- exchange.ResponseFromTCP()
}
