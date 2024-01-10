package tcp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	conn   net.Conn      // conn is the underlying TCP connection.
	reader *bufio.Reader // reader is used for reading data from the connection.
	writer *bufio.Writer // writer is used for writing data to the connection.
	mutex  sync.Mutex    // mutex is a mutex for ensuring thread-safe access to connection-specific operations.
	o      chan *transport.Exchange
	i      chan *transport.Exchange
}

func (c *Connection) response() {
	for {
		exchange := transport.New()
		var length uint32
		if err := binary.Read(c.reader, binary.LittleEndian, &length); err != nil {
			exchange.Response().Body = []byte(fmt.Sprintf("Échec de la lecture de la longueur de la réponse: %v", err))
			exchange.Response().Status = http.StatusInternalServerError
			c.i <- exchange
			continue
		}

		responseBytes := make([]byte, length)
		if _, err := c.reader.Read(responseBytes); err != nil {
			exchange.Response().Body = []byte(fmt.Sprintf("Échec de la lecture de la réponse: %v", err))
			exchange.Response().Status = http.StatusInternalServerError
			c.i <- exchange
			continue
		}

		if err := proto.Unmarshal(responseBytes, exchange.Response()); err != nil {
			logger.Debug(string(responseBytes), err.Error())
			exchange.Response().Body = []byte(fmt.Sprintf("Échec de la désérialisation de la réponse: %v", err))
			exchange.Response().Status = http.StatusInternalServerError
			c.i <- exchange
			continue
		}

		c.i <- exchange
	}
}

func (c *Connection) request() {
	for r := range c.o {
		req := r.Request()
		requestBytes, err := proto.Marshal(req)
		if err != nil {
			r.Response().Body = []byte(fmt.Sprintf("Échec de la sérialisation de la requête: %v", err))
			r.Response().Status = http.StatusInternalServerError
			continue
		}

		if err := binary.Write(c.writer, binary.LittleEndian, uint32(len(requestBytes))); err != nil {
			r.Response().Body = []byte(fmt.Sprintf("Échec de l'écriture de la longueur de la requête: %v", err))
			r.Response().Status = http.StatusInternalServerError
			continue
		}

		if _, err := c.writer.Write(requestBytes); err != nil {
			r.Response().Body = []byte(fmt.Sprintf("Échec de l'écriture de la requête: %v", err))
			r.Response().Status = http.StatusInternalServerError
			continue
		}

		if err := c.writer.Flush(); err != nil {
			r.Response().Body = []byte(fmt.Sprintf("Échec de l'envoi de la requête: %v", err))
			r.Response().Status = http.StatusInternalServerError
			continue
		}
	}
}

/*

	reader := conn.reader
	writer := conn.writer

	req := exchange.Request()
	requestBytes, err := proto.Marshal(req)
	if err != nil {
		exchange.Response().Body = []byte(fmt.Sprintf("Échec de la sérialisation de la requête: %v", err))
		exchange.Response().Status = http.StatusInternalServerError
		return exchange
	}

	if err := binary.Write(writer, binary.LittleEndian, uint32(len(requestBytes))); err != nil {
		exchange.Response().Body = []byte(fmt.Sprintf("Échec de l'écriture de la longueur de la requête: %v", err))
		exchange.Response().Status = http.StatusInternalServerError
		return exchange
	}

	if _, err := writer.Write(requestBytes); err != nil {
		exchange.Response().Body = []byte(fmt.Sprintf("Échec de l'écriture de la requête: %v", err))
		exchange.Response().Status = http.StatusInternalServerError
		return exchange
	}

	if err := writer.Flush(); err != nil {
		exchange.Response().Body = []byte(fmt.Sprintf("Échec de l'envoi de la requête: %v", err))
		exchange.Response().Status = http.StatusInternalServerError
		return exchange
	}
	/*
		var length uint32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			exchange.Response().Body = []byte(fmt.Sprintf("Échec de la lecture de la longueur de la réponse: %v", err))
			exchange.Response().Status = http.StatusInternalServerError
			return exchange
		}

		responseBytes := make([]byte, length)
		if _, err := reader.Read(responseBytes); err != nil {
			exchange.Response().Body = []byte(fmt.Sprintf("Échec de la lecture de la réponse: %v", err))
			exchange.Response().Status = http.StatusInternalServerError
			return exchange
		}

		if err := proto.Unmarshal(responseBytes, exchange.Response()); err != nil {
			exchange.Response().Body = []byte(fmt.Sprintf("Échec de la désérialisation de la réponse: %v", err))
			exchange.Response().Status = http.StatusInternalServerError
			return exchange
		}
*/

// newConnection creates a new instance of the Connection structure
// that encapsulates an underlying TCP connection. It initializes the structure with
// the necessary parameters and attributes for managing the connection.
//
// Parameters:
// - conn: net.Conn The underlying TCP connection object to encapsulate.
//
// Returns:
// - *Connection: A pointer to the newly created Connection instance.
func newConnection(address string, i chan *transport.Exchange) *Connection {
	var conn net.Conn
	var err error

	conn, err = net.Dial("tcp", address)
	for logger.Error(err) {
		conn, err = net.Dial("tcp", address)
	}

	c := &Connection{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		o:      make(chan *transport.Exchange),
		i:      i,
	}

	go c.response()
	go c.request()

	return c
}
