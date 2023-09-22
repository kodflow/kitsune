package router

import (
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

func Resolve(req *transport.Request) *transport.Response {
	res := transport.ResponseFromRequest(req)
	// TODO CONTROLLER

	return res
}
