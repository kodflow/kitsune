package router

import (
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

var Empty = []byte{}

func Resolve(req *transport.Request, res *transport.Response) error {
	res.Id = req.Id
	res.Pid = req.Pid
	return nil
}
