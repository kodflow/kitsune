package tcp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"google.golang.org/protobuf/proto"
)

// Service represents a single service connection.
// It encapsulates the details for a TCP connection, including network I/O objects and a mutex for synchronization.
type Service struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	mutex   sync.Mutex // Mutex for synchronizing network I/O operations
}

func (s *Service) connect() error {
	var err error
	s.conn, err = net.Dial("tcp", s.address)
	if err != nil {
		return err
	}

	s.reader = bufio.NewReader(s.conn)
	s.writer = bufio.NewWriter(s.conn)

	return nil
}

// Send sends a request to the TCP server and waits for a response.
func (s *Service) Send(exchange *transport.Exchange) *transport.Exchange {
	s.mutex.Lock() // Verrouillage avant l'accès partagé
	defer s.mutex.Unlock()

	req := exchange.Request()
	res := exchange.Response()

	requestBytes, err := proto.Marshal(req)
	if err != nil {
		res.Body = []byte(fmt.Sprintf("failed to marshal request: %v", err))
		res.Status = http.StatusInternalServerError
	}

	if err := binary.Write(s.writer, binary.LittleEndian, uint32(len(requestBytes))); err != nil {
		res.Body = []byte(fmt.Sprintf("failed to write request length: %v", err))
		res.Status = http.StatusInternalServerError
	}

	if _, err := s.writer.Write(requestBytes); err != nil {
		res.Body = []byte(fmt.Sprintf("failed to write request: %v", err))
		res.Status = http.StatusInternalServerError
	}

	if err := s.writer.Flush(); err != nil {
		res.Body = []byte(fmt.Sprintf("failed to flush request: %v", err))
		res.Status = http.StatusInternalServerError
	}

	var length uint32
	if err := binary.Read(s.reader, binary.LittleEndian, &length); err != nil {
		res.Body = []byte(fmt.Sprintf("failed to read response length: %v", err))
		res.Status = http.StatusInternalServerError
	}

	responseBytes := make([]byte, length)
	if _, err := s.reader.Read(responseBytes); err != nil {
		res.Body = []byte(fmt.Sprintf("failed to read response: %v", err))
		res.Status = http.StatusInternalServerError
	}

	if err := proto.Unmarshal(responseBytes, res); err != nil {
		res.Body = []byte(fmt.Sprintf("failed to unmarshal response: %v", err))
		res.Status = http.StatusInternalServerError
	}

	return exchange
}

// Close closes the TCP connection of the service.
// It is responsible for releasing the network resources associated with the service.
//
// Returns:
// - error: Error, if any occurred during the closure of the connection.
func (s *Service) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
