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
	ip := "127.0.0.1:" + generateRandomNumbers()
	t.Run("New", func(t *testing.T) {
		server := setupServer(ip)
		assert.NotNil(t, server)
		assert.Equal(t, ip, server.Address)
	})

	t.Run("StartAndStop", func(t *testing.T) {
		server := setupServer(ip)
		assert.NotNil(t, server)

		assert.Error(t, server.Stop())
		assert.NoError(t, server.Start())
		assert.Error(t, server.Start())
		assert.NoError(t, server.Stop())
		assert.Error(t, server.Stop())
	})

	t.Run("Start:Successfully", func(t *testing.T) {
		server := setupServer(ip)
		err := server.Start()
		assert.Nil(t, err)

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		server.Stop()
	})

	t.Run("Start:Failure(already started)", func(t *testing.T) {
		server := setupServer(ip)
		server.Start()

		// Attempting to start again
		err := server.Start()
		assert.NotNil(t, err)
		assert.Equal(t, "server already started", err.Error())

		server.Stop()
	})

	t.Run("Stop:Successfully", func(t *testing.T) {
		server := setupServer(ip)
		server.Start()

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		err := server.Stop()
		assert.Nil(t, err)
	})

	t.Run("Stop:Failure(already stopped)", func(t *testing.T) {
		server := setupServer(ip)
		err := server.Stop()
		assert.NotNil(t, err)
		assert.Equal(t, "server is not active", err.Error())
	})
}
