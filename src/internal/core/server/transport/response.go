package transport

import "google.golang.org/protobuf/proto"

func ResponseFromBytes(b []byte) *Response {
	res := &Response{}
	proto.Unmarshal(b, res)

	return res
}

func ResponseToBytes(res *Response) ([]byte, error) {
	return proto.Marshal(res)
}

func ResponseFromRequest(req *Request) *Response {
	return &Response{Id: req.Id}
}
