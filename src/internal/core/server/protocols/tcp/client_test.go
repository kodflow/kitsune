package tcp

import (
	"sync"
	"testing"
	"time"

	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

func TestTCPClient(t *testing.T) {
	config.DEFAULT_LOG_LEVEL = levels.DEBUG
	server := setupServer("127.0.0.1:7778")
	server.Start()
	defer server.Stop()

	client := NewClient()
	service, err := client.Connect(server.Address, 20)
	assert.Nil(t, err)
	service2, err := client.Connect(server.Address, 20)
	assert.Nil(t, err)
	assert.Equal(t, service, service2)

	i := transport.New()
	o := service.Send(i)

	o.Wait()

	assert.Equal(t, i.Request().Id, o.Request().Id)
	assert.Equal(t, o.Request().Id, o.Response().Id)
	assert.Nil(t, err)

	requestMax := 1000

	var mu sync.Mutex

	tick := 0
	v := 0
	ticker := time.NewTicker(time.Second)
	I := 0
	O := 0

	go func() {

		logger.Warn("Start")
		for v = 0; v < requestMax; v++ {
			go func() {
				i := transport.New()
				assert.Nil(t, i.Response())
				mu.Lock()
				I++
				mu.Unlock()
				o := service.Send(i)
				o.Wait()
				mu.Lock()
				O++
				mu.Unlock()
				assert.Equal(t, i.Request().Id, o.Request().Id)
				assert.Equal(t, i.Request().Id, o.Response().Id)
			}()
		}
		logger.Warn("Stop")
	}()

	for range ticker.C {
		tick++
		if O == requestMax {
			ticker.Stop()
			client.Close()
			logger.Success("Finished")
			break
		}
	}
}
