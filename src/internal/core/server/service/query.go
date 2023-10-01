package service

import (
	"github.com/google/uuid"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

type Query struct {
	Answer  bool
	Service string
	Req     *transport.Request
}

func query(service string, answer bool) *Query {
	return &Query{Service: service, Req: request(), Answer: answer}
}

func request() *transport.Request {
	v4, err := uuid.NewRandom()

	if err != nil {
		return nil
	}

	return &transport.Request{Id: v4.String()}
}
