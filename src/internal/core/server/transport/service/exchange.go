// Package service provides the functionality to interact with a service over a network.
package service

import (
	"github.com/google/uuid"
	"github.com/kodmain/kitsune/src/internal/core/server/transport/proto/generated"
)

// Exchange struct represents a service message.
type Exchange struct {
	Service string             // The service to which this message is directed.
	Answer  bool               // Indicates if this message expects a response.
	Req     *generated.Request // The transport request associated with this message.
}

// NewExchange creates a new Exchange instance.
//
// Parameters:
// - service: The service to which the exchange is directed.
// - answer: Whether the exchange expects a response or not.
//
// Returns:
// - exchange: A pointer to the newly created exchange object.
func NewExchange(service string, answer bool) *Exchange {
	v4, err := uuid.NewRandom()

	if err != nil {
		return nil
	}

	return &Exchange{
		Service: service,
		Req: &generated.Request{
			Endpoint: service,
			Id:       v4.String(),
		},
		Answer: answer,
	}
}
