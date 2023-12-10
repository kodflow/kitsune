package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

// HTTPClient represents an HTTP client.
// It provides methods to send HTTP requests and receive responses.
type HTTPClient struct {
	client    *http.Client
	transport *http.Transport
}

// NewHTTPClient initializes a new HTTPClient and returns its pointer.
//
// Returns:
//
//	*HTTPClient: A pointer to the newly created HTTPClient
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: config.DEFAULT_TIMEOUT * time.Second,
		},
	}
}

// SendRequest sends an HTTP request and returns a response or an error.
//
// Parameters:
//
//	method: The HTTP method (GET, POST, etc.)
//	url: The URL for the request
//	body: The request body (for POST, PUT, etc.)
//
// Returns:
//
//	[]byte: The response body
//	error: An error if one occurred
func (c *HTTPClient) Send(req *transport.Request) *transport.Response {
	res := &transport.Response{
		Status:  500,
		Headers: map[string]string{},
	}

	httpRequest, err := http.NewRequest(req.Method, req.Endpoint, bytes.NewReader(req.Body))

	if err != nil {
		return res
	}

	for k, v := range req.Headers {
		httpRequest.Header.Set(k, v)
	}

	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		return res
	}

	defer httpResponse.Body.Close()

	res.Status = uint32(httpResponse.StatusCode)
	for header := range httpResponse.Header {
		res.Headers[header] = httpResponse.Header.Get(header)
	}

	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return res
	}

	res.Body = data

	return res
}
