package socket_test

import (
	"testing"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/socket"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/stretchr/testify/assert"
)

func TestPromise(t *testing.T) {
	server := socket.NewServer("localhost:8080")
	server.Start()
	defer server.Stop()

	client := socket.NewClient("localhost:8080")

	t.Run("Promise", func(t *testing.T) {
		client.Connect()
		req := transport.CreateRequest()
		promise, _ := client.Send(req)
		res := promise.Wait()
		assert.NotNil(t, res, "Expected response")
		assert.Equal(t, res.Id, req.Id)
	})

	t.Run("NoPromise", func(t *testing.T) {
		client.Connect()
		req := transport.CreateRequestOnly()
		promise, _ := client.Send(req)
		res := promise.Wait()
		assert.Nil(t, res, "Expected no response")
	})
}
