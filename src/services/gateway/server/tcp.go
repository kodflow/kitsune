package server

import (
	"fmt"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/tcp"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

var Handler *daemon.Handler = &daemon.Handler{
	Name: "process manager",
	Call: func() error {
		s := tcp.NewServer("0.0.0.0:9999")
		fmt.Println("audit on", "0.0.0.0:9999")
		return s.Start()
	},
}
