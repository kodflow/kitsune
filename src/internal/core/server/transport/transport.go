package transport

import "github.com/google/uuid"

var Empty = []byte{}

func New() (*Request, *Response) {
	id, _ := uuid.NewRandom()

	req := &Request{
		Id: id.String(),
	}

	res := &Response{
		Status: 204,
		Id:     req.Id,
	}

	return req, res
}
