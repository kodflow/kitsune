package server

import (
	"github.com/kodmain/kitsune/src/internal/core/server/protocols/socket"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

var Handler *daemon.Handler = &daemon.Handler{
	Name: "process manager",
	Call: func() error {
		s := socket.NewServer("localhost:9999")

		return s.Start()
	},
}
