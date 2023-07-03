package query

import "github.com/google/uuid"

type Message struct {
	ID       string
	Request  *Request
	Response *Response
}

func NewMessage(req *Request) *Message {
	return &Message{
		ID:      uuid.NewString(),
		Request: req,
	}
}
