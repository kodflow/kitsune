package tcp

import (
	"net"
	"testing"

	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	// Create a mock channel for the response
	responseChan := make(chan *generated.Response)

	// Start a mock TCP server
	listener, err := net.Listen("tcp", "localhost:0")
	assert.Nil(t, err)
	// Create a new connection
	address := listener.Addr().String()
	conn := newConnection(address, responseChan)

	// Perform assertions on the connection object
	assert.NotNil(t, conn)
	assert.NotNil(t, conn.net)
	assert.NotNil(t, conn.reader)
	assert.NotNil(t, conn.writer)
	assert.NotNil(t, conn.o)
	assert.Equal(t, responseChan, conn.i)

	// Clean up resources
	close(responseChan)
}
