package main

import (
	"github.com/kodmain/kitsune/src/internal/core/server/protocols/http"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

func main() {
	daemon.New().Start(&daemon.Handler{
		Name: "HTTP Server",
		Call: func() error {
			server := http.NewServer(&http.ServerCfg{
				HTTP:  80,
				HTTPS: 443,
				//DOMAIN: "aube.io",
				//SUBS:   []string{"home"},
			})

			return server.Start()
		},
	})
}
