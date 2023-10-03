// Package service manages service-level functionalities.
package service

import (
	"github.com/google/uuid"                                        // Importing package for generating UUIDs.
	"github.com/kodmain/kitsune/src/internal/core/server/transport" // Importing package for transport-level logic.
)

// Query struct models a service query.
type Query struct {
	Answer  bool               // Indicates if this query expects an answer.
	Service string             // The service to which this query is directed.
	Req     *transport.Request // The actual transport request associated with this query.
}

// query creates a new Query instance.
// service: The service to which the query is directed.
// answer: Whether the query expects a response or not.
// returns a pointer to the newly created Query object.
func query(service string, answer bool) *Query {
	return &Query{Service: service, Req: request(), Answer: answer}
}

// request creates a new transport.Request with a random UUID.
// returns a pointer to a new Request object, or nil if an error occurs.
func request() *transport.Request {
	v4, err := uuid.NewRandom() // Generate a random UUID.

	// If generating a UUID fails, return nil.
	if err != nil {
		return nil
	}

	// Create and return a new Request with the generated UUID.
	return &transport.Request{Id: v4.String()}
}
