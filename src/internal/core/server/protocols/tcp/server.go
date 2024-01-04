package tcp

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/kodflow/kitsune/src/internal/core/server/router"
	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

// Server represents a TCP server and contains information about the address it listens on
// and the underlying network listener.
type Server struct {
	Address  string       // Address to listen on
	listener net.Listener // TCP Listener object
	router   *router.Router
}

// NewServer creates a new Server instance with the specified listening address.
func NewServer(address string) *Server {
	return &Server{
		Address: address,
		router:  router.MakeRouter(),
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
	if logger.Error(err) {
		return err
	}

	go s.accepLoop()

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

	logger.Info("server stop on " + s.Address)
	return err
}

// accepLoop continuously accepts incoming connections.
// It listens for incoming client connections and handles them asynchronously by calling 'handleConnection'.
func (s *Server) accepLoop() {
	for {
		if s.listener == nil {
			break
		}

		conn, err := s.listener.Accept() // Accept incoming connections.
		if err != nil {
			break // Exit the loop if there is an error accepting a connection.
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
	defer conn.Close()              // Close the connection when the function exits.
	reader := bufio.NewReader(conn) // Create a buffered reader for reading data from the connection.
	writer := bufio.NewWriter(conn) // Create a buffered writer for writing data to the connection.
	for {
		data, err := s.handleRequest(reader) // Read and process incoming requests.
		if logger.Error(err) {
			break // Exit the loop if there is an error.
		}

		s.sendResponse(writer, s.TCPHandler(data)) // Send a response to the client based on the processed data.
	}
}

// handleRequest reads and processes incoming requests.
//
// Parameters:
// - reader: *bufio.Reader - The reader used to read incoming data.
//
// Returns:
// - []byte: The data read from the reader.
// - error: An error if there was an issue reading the request.
func (s *Server) handleRequest(reader *bufio.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err := io.ReadFull(reader, data)
	return data, err
}

// sendResponse sends a response to the client.
//
// Parameters:
// - writer: *bufio.Writer - The writer used to send the response.
// - res: []byte - The response data to be sent.
func (s *Server) sendResponse(writer *bufio.Writer, res []byte) {
	if len(res) > 0 {
		binary.Write(writer, binary.LittleEndian, uint32(len(res)))
		writer.Write(res)
		writer.Flush()
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
func (s *Server) TCPHandler(b []byte) []byte {
	// Initialize a new transport request and response
	req, res := transport.New()

	// Unmarshal the input byte array into the request struct
	err := proto.Unmarshal(b, req)

	// Set the process ID in the response to match the request
	res.Pid = req.Pid

	// Return an empty response if there's an error in unmarshalling
	if err != nil {
		return transport.Empty
	}

	// Resolve the request using the router and update the response
	err = s.router.Resolve(req, res)

	// Return an empty response if there's an error in processing the request
	if err != nil {
		return transport.Empty
	}

	// Marshal the response back into a byte array
	b, err = proto.Marshal(res)

	// Return an empty response if there's an error in marshalling
	if err != nil {
		return transport.Empty
	}

	// Return the processed response
	return b
}
