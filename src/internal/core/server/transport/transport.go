package transport

import "github.com/google/uuid"

// Empty represents an empty byte slice.
// It can be used as a placeholder or default value where a byte slice is required but no data is provided.
var Empty = []byte{}

// New creates and returns a new Request and Response pair with a unique ID.
// It utilizes the UUID library to generate a unique identifier for each request-response pair.
// The Response is initialized with a default status of 204 (No Content) and the same ID as the Request.
//
// Returns:
// - *Request: A pointer to a newly created Request object with a unique ID.
// - *Response: A pointer to a newly created Response object with a status of 204 and the same ID as the Request.
func New() (*Request, *Response) {
	id, _ := uuid.NewRandom() // Generate a new random UUID.

	req := &Request{
		Id: id.String(), // Set the Request ID to the UUID string.
	}

	res := &Response{
		Status: 204,    // Initialize with a 204 No Content status.
		Id:     req.Id, // Set the Response ID to the same ID as the Request.
	}

	return req, res
}
