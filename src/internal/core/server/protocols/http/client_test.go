package http_test

import (
	"testing"

	"github.com/kodflow/kitsune/src/internal/core/server/protocols/http"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	// Création d'un serveur HTTP de test
	server := setupHTTPServer(7777, 7778)
	server.Start()
	defer server.Stop()

	// Initialisation du client HTTP
	client := http.NewHTTPClient()

	req := &generated.Request{
		Method:   "GET",
		Endpoint: "http://127.0.0.1:7777/",
	}

	t.Run("GET Request", func(t *testing.T) {
		res := client.Send(req)
		assert.Equal(t, uint32(204), res.Status)
	})

	req = &generated.Request{
		Method:   "GET",
		Endpoint: "https://127.0.0.1:7778/",
	}

	t.Run("GET Request", func(t *testing.T) {
		res := client.Send(req)
		assert.Equal(t, uint32(204), res.Status)
	})
}
