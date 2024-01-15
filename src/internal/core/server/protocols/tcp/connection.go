package tcp

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	close  bool
	net    net.Conn      // net is the underlying TCP connection.
	reader *bufio.Reader // reader is used for reading data from the connection.
	writer *bufio.Writer // writer is used for writing data to the connection.
	mutex  sync.Mutex    // mutex is a mutex for ensuring thread-safe access to connection-specific operations.
	o      chan *generated.Request
	i      chan *generated.Response
}

func (c *Connection) response() {
	for {
		if c.net == nil {
			break
		}
		var length uint32

		err := binary.Read(c.net, binary.LittleEndian, &length)
		if c.close || logger.Error(err) {
			continue
		}

		// Lire directement les données de réponse depuis net.Conn
		responseBytes := make([]byte, length)
		_, err = io.ReadFull(c.net, responseBytes)

		if logger.Error(err) {
			continue
		}

		// Désérialisation de la réponse
		var res *generated.Response = transport.NewReponse()
		if err := proto.Unmarshal(responseBytes, res); logger.Error(err) {
			continue
		}

		c.i <- res
	}
}

func (c *Connection) request() {
	for req := range c.o {
		requestBytes, err := proto.Marshal(req)
		if logger.Error(err) {
			continue
		}

		// Préparer la longueur de la requête en tant que préfixe
		lengthBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lengthBytes, uint32(len(requestBytes)))

		// Écrire la longueur suivie de la requête directement sur net.Conn
		if _, err = c.net.Write(lengthBytes); logger.Error(err) {
			continue
		}

		if _, err = c.net.Write(requestBytes); logger.Error(err) {
			continue
		}
	}
}

// newConnection creates a new instance of the Connection structure
// that encapsulates an underlying TCP connection. It initializes the structure with
// the necessary parameters and attributes for managing the connection.
//
// Parameters:
// - conn: net.Conn The underlying TCP connection object to encapsulate.
//
// Returns:
// - *Connection: A pointer to the newly created Connection instance.
func newConnection(address string, i chan *generated.Response) *Connection {
	var conn net.Conn
	var err error

	conn, err = net.Dial("tcp", address)
	for logger.Error(err) {
		time.Sleep(time.Second)
		conn, err = net.Dial("tcp", address)
	}

	c := &Connection{
		net:    conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		o:      make(chan *generated.Request),
		i:      i,
	}

	go c.response()
	go c.request()

	return c
}
