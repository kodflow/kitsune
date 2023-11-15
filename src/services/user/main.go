package main

import (
	"github.com/kodmain/kitsune/src/internal/core/server/protocols/tcp"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
	user "github.com/kodmain/kitsune/src/services/user/api"
)

func main() {
	daemon.New().Start(&daemon.Handler{
		Name: "TCP Server",
		Call: func() error {
			server := tcp.NewServerV2(":9000")
			server.Register(user.V1) // API V1
			server.Register(user.V2) // API V2
			return server.Start()
		},
	})
}
