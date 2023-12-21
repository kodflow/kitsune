package http_test

import (
	"testing"
	"time"

	"github.com/kodflow/kitsune/src/internal/core/server/protocols/http"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

func setupHTTPServer(httpPort, httpsPort string) *http.Server {
	cfg := &http.ServerCfg{
		DOMAIN: "127.0.0.1",
		SUBS:   []string{},
		HTTP:   httpPort,
		HTTPS:  httpsPort,
	}
	return http.NewServer(cfg)
}

func TestHTTPServer(t *testing.T) {
	p1, p2 := generateTwoDistinctRandomNumbers()
	// Création d'un serveur HTTP de test
	logger.SetLevel(levels.OFF)
	t.Run("NewHTTPServer", func(t *testing.T) {
		server := setupHTTPServer(p1, "")
		assert.NotNil(t, server)
	})

	t.Run("StartStandardServer:Successfully", func(t *testing.T) {
		server := setupHTTPServer(p1, "")
		err := server.Start()
		assert.Nil(t, err)

		time.Sleep(100 * time.Millisecond)

		server.Stop()
	})

	t.Run("StartStandardServer:Failure(already started)", func(t *testing.T) {
		server := setupHTTPServer(p1, "")
		server.Start()

		err := server.Start()
		assert.NotNil(t, err)

		server.Stop()
	})

	t.Run("StartSecureServer:Successfully", func(t *testing.T) {
		server := setupHTTPServer("", p2)
		err := server.Start()
		assert.Nil(t, err)

		time.Sleep(100 * time.Millisecond)

		server.Stop()
	})

	t.Run("StartSecureServer:Failure(already started)", func(t *testing.T) {
		server := setupHTTPServer("", p2)
		server.Start()

		err := server.Start()
		assert.NotNil(t, err)

		server.Stop()
	})
}
