package promise_test

import (
	"testing"
	"time"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/kodmain/kitsune/src/internal/core/server/transport/promise"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	config.DEFAULT_TIMEOUT = 1
	callbackCalled := false
	callback := func(responses ...*transport.Response) {
		callbackCalled = true
	}

	p, err := promise.Create(callback)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	// Wait for the timeout duration to ensure the promise is closed.
	time.Sleep(time.Second*config.DEFAULT_TIMEOUT + time.Second)

	assert.True(t, p.Closed)
	assert.True(t, callbackCalled)
}
