package transport

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func RequestFromBytes(b []byte) *Request {
	req := &Request{}
	proto.Unmarshal(b, req)

	return req
}

func CreateRequestOnly() *Request {
	v4, err := uuid.NewRandom()

	if err != nil {
		return nil
	}

	return &Request{Id: v4.String()}
}

func CreateRequest() *Request {
	v4, err := uuid.NewRandom()

	if err != nil {
		return nil
	}

	return &Request{Id: v4.String(), Answer: true}
}
