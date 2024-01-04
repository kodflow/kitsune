package main

import (
	"github.com/kodflow/kitsune/src/internal/core/server/protocols/tcp"
	"github.com/kodflow/kitsune/src/internal/kernel/daemon"
)

func main() {
	daemon.New().Start(&daemon.Handler{
		Name: "TCP Server",
		Call: func() error {
			server := tcp.NewServer(":9999")
			/*
				server.Register(user.V1) // API V1
				server.Register(user.V2) // API V2
			*/
			return server.Start()
		},
	})
}
