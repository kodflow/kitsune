package http_test

import (
	"testing"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/http"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	// Création d'un serveur HTTP de test
	server := setupHTTPServer(80, 443)
	server.Start()
	defer server.Stop()

	// Initialisation du client HTTP
	client := http.NewHTTPClient()

	req := &transport.Request{
		Method:   "GET",
		Endpoint: "http://localhost/",
	}

	t.Run("GET Request", func(t *testing.T) {
		res := client.Send(req)
		assert.Equal(t, uint32(204), res.Status)
	})

	req = &transport.Request{
		Method:   "GET",
		Endpoint: "https://localhost/",
	}

	t.Run("GET Request", func(t *testing.T) {
		res := client.Send(req)
		assert.Equal(t, uint32(204), res.Status)
	})
}
