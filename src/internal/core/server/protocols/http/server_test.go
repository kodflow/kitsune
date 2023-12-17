package http_test

import (
	"testing"
	"time"

	"github.com/kodflow/kitsune/src/internal/core/server/protocols/http"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

func setupHTTPServer(httpPort, httpsPort int) *http.Server {
	cfg := &http.ServerCfg{
		DOMAIN: "127.0.0.1",
		SUBS:   []string{},
		HTTP:   httpPort,
		HTTPS:  httpsPort,
	}
	return http.NewServer(cfg)
}

func TestHTTPServer(t *testing.T) {
	logger.SetLevel(levels.OFF)
	t.Run("NewHTTPServer", func(t *testing.T) {
		server := setupHTTPServer(8080, 0)
		assert.NotNil(t, server)
	})

	t.Run("StartStandardServer:Successfully", func(t *testing.T) {
		server := setupHTTPServer(8080, 0)
		err := server.Start()
		assert.Nil(t, err)

		time.Sleep(100 * time.Millisecond)

		server.Stop()
	})

	t.Run("StartStandardServer:Failure(already started)", func(t *testing.T) {
		server := setupHTTPServer(8080, 0)
		server.Start()

		err := server.Start()
		assert.NotNil(t, err)

		server.Stop()
	})

	t.Run("StartSecureServer:Successfully", func(t *testing.T) {
		server := setupHTTPServer(0, 8443)
		err := server.Start()
		assert.Nil(t, err)

		time.Sleep(100 * time.Millisecond)

		server.Stop()
	})

}
