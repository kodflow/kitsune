// Package service provides the functionality to interact with a service over a network.
package service

import (
	"github.com/google/uuid"                                        // Importing package for generating UUIDs.
	"github.com/kodmain/kitsune/src/internal/core/server/transport" // Importing package for transport-level logic.
)

// Exchange struct models a service message.
type Exchange struct {
	service string             // The service to which this message is directed.
	answer  bool               // Indicates if this message expects an answer.
	req     *transport.Request // The actual transport request associated with this message.
}

func (m *Exchange) ServiceName() string {
	return m.service
}

func (m *Exchange) Answer() bool {
	return m.answer
}

func (m *Exchange) Request() *transport.Request {
	return m.req
}

// NewExchange creates a new Exchange instance.
// service: The service to which the exchange is directed.
// answer: Whether the exchange expects a response or not.
// returns a pointer to the newly created exchange object.
func NewExchange(service string, answer bool) *Exchange {
	v4, err := uuid.NewRandom()

	if err != nil {
		return nil
	}

	return &Exchange{
		service: service,
		req:     &transport.Request{Id: v4.String()},
		answer:  answer,
	}
}
