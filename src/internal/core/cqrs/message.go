// Package service manages service-level functionalities.
package cqrs

import (
	"github.com/google/uuid"                                        // Importing package for generating UUIDs.
	"github.com/kodmain/kitsune/src/internal/core/server/transport" // Importing package for transport-level logic.
)

type CQRS bool

const QRY CQRS = false
const CMD CQRS = true

// Message struct models a service message.
type Message struct {
	cqrs    CQRS
	service string             // The service to which this message is directed.
	answer  bool               // Indicates if this message expects an answer.
	req     *transport.Request // The actual transport request associated with this message.
}

func (m *Message) ServiceName() string {
	return m.service
}

func (m *Message) Answer() bool {
	return m.answer
}

func (m *Message) Request() *transport.Request {
	return m.req
}

// NewMessage creates a new Message instance.
// service: The service to which the message is directed.
// answer: Whether the message expects a response or not.
// returns a pointer to the newly created Message object.
func NewMessage(service string, answer bool, cqrs CQRS) *Message {
	v4, err := uuid.NewRandom()

	if err != nil {
		return nil
	}

	return &Message{
		service: service,
		req:     &transport.Request{Id: v4.String()},
		answer:  answer,
		cqrs:    cqrs,
	}
}
