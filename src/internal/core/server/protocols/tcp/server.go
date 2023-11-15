// Package tcp provides functionalities for a TCP server.
package tcp

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/kodmain/kitsune/src/internal/core/server/api"
	"github.com/kodmain/kitsune/src/internal/core/server/router"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

//https://github.com/douglasmakey/socket-sharding/blob/master/cmd/http-example/main.go

// Server represents a TCP server and contains information about the address it listens on
// and the underlying network listener.
type Server struct {
	Address  string       // Address to listen on
	listener net.Listener // TCP Listener object
}

// NewServer creates a new Server instance with the specified listening address.
func NewServer(address string) *Server {
	return &Server{
		Address: address,
	}
}

// Register is a method for registering API handlers with the server.
func (s *Server) Register(api api.APInterface) {

}

// Start starts the TCP server, allowing it to accept incoming connections.
func (s *Server) Start() error {
	if s.listener != nil {
		return errors.New("server already started")
	}

	var err error
	s.listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	go s.accepLoop()

	logger.Info("server start on " + s.Address + " with pid:" + strconv.Itoa(os.Getpid()))

	return nil
}

// Stop stops the TCP server.
func (s *Server) Stop() error {
	if s.listener == nil {
		return errors.New("server is not active")
	}

	err := s.listener.Close()
	s.listener = nil

	logger.Info("server stop on " + s.Address)
	return err
}

// acceptLoop continuously accepts incoming connections.
func (s *Server) accepLoop() {
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
}

// handleConnection handles incoming client connections.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		data, err := s.handleRequest(reader)
		if err != nil {
			break
		}

		s.sendResponse(writer, s.handler(data))
	}
}

// handleRequest reads and processes incoming requests.
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
func (s *Server) sendResponse(writer *bufio.Writer, res []byte) {
	if len(res) > 0 {
		binary.Write(writer, binary.LittleEndian, uint32(len(res)))
		writer.Write(res)
		writer.Flush()
	}
}

// handler processes incoming data and generates a response.
func (s *Server) handler(b []byte) []byte {
	req := &transport.Request{}
	res := &transport.Response{}
	err := proto.Unmarshal(b, req)
	if err != nil {
		return router.Empty
	}

	err = router.Resolve(req, res)
	if err != nil {
		return router.Empty
	}

	b, err = proto.Marshal(res)
	if err != nil {
		return router.Empty
	}

	return b
}
