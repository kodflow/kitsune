package main

import (
	"github.com/kodmain/kitsune/src/internal/core/server/protocols/http"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

func main() {
	daemon.New().Start(&daemon.Handler{
		Name: "HTTP Server",
		Call: func() error {
			return http.NewServer("gateway.kitsune.local:8888").Start()
		},
	})
}
