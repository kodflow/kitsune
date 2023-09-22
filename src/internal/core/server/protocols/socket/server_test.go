package socket_test

import (
	"testing"
	"time"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/socket"
	"github.com/stretchr/testify/assert"
)

func setup() *socket.Server {
	return socket.NewServer("localhost:8080")
}

func TestServer(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		server := setup()
		assert.NotNil(t, server)
		assert.Equal(t, "localhost:8080", server.Address)
	})

	t.Run("Start:Successfully", func(t *testing.T) {
		server := setup()
		err := server.Start()
		assert.Nil(t, err)

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		server.Stop()
	})

	t.Run("Start:Failure(already started)", func(t *testing.T) {
		server := setup()
		server.Start()

		// Attempting to start again
		err := server.Start()
		assert.NotNil(t, err)
		assert.Equal(t, "server already started", err.Error())

		server.Stop()
	})

	t.Run("Stop:Successfully", func(t *testing.T) {
		server := setup()
		server.Start()

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		err := server.Stop()
		assert.Nil(t, err)
	})

	t.Run("Stop:Failure(already stopped)", func(t *testing.T) {
		server := setup()
		err := server.Stop()
		assert.NotNil(t, err)
		assert.Equal(t, "server is not active", err.Error())
	})

}