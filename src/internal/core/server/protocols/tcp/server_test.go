package tcp

import (
	"testing"
	"time"

	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

func setupServer(address string) *Server {
	return NewServer(address)
}

func TestServer(t *testing.T) {
	logger.SetLevel(levels.OFF)
	t.Run("New", func(t *testing.T) {
		server := setupServer("127.0.0.1:7777")
		assert.NotNil(t, server)
		assert.Equal(t, "127.0.0.1:7777", server.Address)
	})

	t.Run("Start:Successfully", func(t *testing.T) {
		server := setupServer("127.0.0.1:7777")
		err := server.Start()
		assert.Nil(t, err)

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		server.Stop()
	})

	t.Run("Start:Failure(already started)", func(t *testing.T) {
		server := setupServer("127.0.0.1:7777")
		server.Start()

		// Attempting to start again
		err := server.Start()
		assert.NotNil(t, err)
		assert.Equal(t, "server already started", err.Error())

		server.Stop()
	})

	t.Run("Stop:Successfully", func(t *testing.T) {
		server := setupServer("127.0.0.1:7777")
		server.Start()

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		err := server.Stop()
		assert.Nil(t, err)
	})

	t.Run("Stop:Failure(already stopped)", func(t *testing.T) {
		server := setupServer("127.0.0.1:7777")
		err := server.Stop()
		assert.NotNil(t, err)
		assert.Equal(t, "server is not active", err.Error())
	})
}
