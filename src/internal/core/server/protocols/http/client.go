package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
)

// HTTPClient represents an HTTP client.
// It encapsulates the functionality for sending HTTP requests and receiving responses,
// handling underlying client and transport configurations.
type HTTPClient struct {
	client    *http.Client    // The underlying HTTP client used for sending requests.
	transport *http.Transport // The transport configuration for the HTTP client.
}

// NewHTTPClient initializes a new HTTPClient and returns its pointer.
// This function sets up the HTTP client with a custom transport configuration,
// including TLS settings and default timeout.
//
// Returns:
// - *HTTPClient: A pointer to the newly created HTTPClient.
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

// Send sends an HTTP request and returns the HTTP response.
// This method constructs and sends an HTTP request based on the provided transport request,
// handling headers, method, endpoint, and body. It also processes the received HTTP response,
// extracting status, headers, and body.
//
// Parameters:
// - req: *generated.Request The HTTP request to be sent.
//
// Returns:
// - *generated.Response: The HTTP response received.
func (c *HTTPClient) Send(req *generated.Request) *generated.Response {
	// Create a default response with a 500 status code and an empty header.
	res := &generated.Response{
		Status:  500,
		Headers: map[string]*generated.Header{},
	}

	// Create an HTTP request based on the input request.
	httpRequest, err := http.NewRequest(req.Method, req.Endpoint, bytes.NewReader(req.Body))

	if err != nil {
		// If there's an error creating the HTTP request, return the default response.
		return res
	}

	// Set headers for the HTTP request.
	for k, v := range req.Headers {
		for _, h := range v.GetItems() {
			httpRequest.Header.Add(k, h)
		}
	}

	// Send the HTTP request and receive the HTTP response.
	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		// If there's an error sending the request, return the default response.
		return res
	}

	defer httpResponse.Body.Close()

	// Populate the response with the status code from the HTTP response.
	res.Status = uint32(httpResponse.StatusCode)

	// Copy headers from the HTTP response to the response object.

	for header := range httpResponse.Header {
		res.Headers[header] = &generated.Header{Items: httpResponse.Header.Values(header)}
	}

	// Read the response body and store it in the response object.
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		// If there's an error reading the response body, return the default response.
		return res
	}

	res.Body = data

	// Return the populated response object.
	return res
}
