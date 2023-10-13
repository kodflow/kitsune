// Package service provides the functionality to interact with a service over a network.
package service

import (
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"google.golang.org/protobuf/proto"
)

var empty = []byte{}

func Handler(b []byte) []byte {
	req := &transport.Request{}
	res := &transport.Response{}
	err := proto.Unmarshal(b, req)
	if err != nil {
		return empty
	}

	err = Resolve(req, res)
	if err != nil { // todo what to do when handler return err 500 ?
		return empty
	}

	b, err = proto.Marshal(res)
	if err != nil {
		return empty
	}

	return b
}

func Resolve(req *transport.Request, res *transport.Response) error {
	res.Id = req.Id
	res.Pid = req.Pid
	return nil
}
