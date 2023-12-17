package transport_test

import (
	"testing"

	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/stretchr/testify/assert"
)

func TestRequestInitialization(t *testing.T) {
	req := &generated.Request{
		Id:       "test_id",
		Pid:      "test_pid",
		Method:   "GET",
		Endpoint: "/test",
		Body:     []byte("test_body"),
		Headers:  map[string]string{"Content-Type": "application/json"},
	}

	assert.Equal(t, "test_id", req.Id)
	assert.Equal(t, "test_pid", req.Pid)
}

func TestResponseInitialization(t *testing.T) {
	res := &generated.Response{
		Status:  200,
		Id:      "test_id",
		Pid:     "test_pid",
		Body:    []byte("response_body"),
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	assert.Equal(t, uint32(200), res.Status)
}

func TestNewFunction(t *testing.T) {
	req, res := transport.New()

	assert.NotEmpty(t, req.Id)
	assert.Equal(t, req.Id, res.Id)
	assert.Equal(t, uint32(204), res.Status)
}
