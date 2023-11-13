package tcp_test

import (
	"testing"
	"time"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/tcp"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

func setupNew() *tcp.ServerV2 {
	return tcp.NewServerV2("localhost:9000")
}

func TestServerV2(t *testing.T) {
	logger.SetLevel(levels.OFF)
	t.Run("New", func(t *testing.T) {
		server := setupNew()
		assert.NotNil(t, server)
		assert.Equal(t, "localhost:9000", server.Address)
	})

	t.Run("Start:Successfully", func(t *testing.T) {
		server := setupNew()
		err := server.Start()
		assert.Nil(t, err)

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		err = server.Stop()
		assert.Nil(t, err)
	})

	t.Run("Start:Failure(already started)", func(t *testing.T) {
		server := setupNew()
		server.Start()

		// Attempting to start again
		err := server.Start()
		assert.NotNil(t, err)
		assert.Equal(t, "server already started", err.Error())

		server.Stop()
	})

	t.Run("Stop:Successfully", func(t *testing.T) {
		server := setupNew()
		server.Start()

		// Allow the server some time to start
		time.Sleep(100 * time.Millisecond)

		err := server.Stop()
		assert.Nil(t, err)
	})

	t.Run("Stop:Failure(already stopped)", func(t *testing.T) {
		server := setupNew()
		err := server.Stop()
		assert.NotNil(t, err)
		assert.Equal(t, "server is not active", err.Error())
	})

}
