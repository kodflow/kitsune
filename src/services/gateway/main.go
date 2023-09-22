package main

import (
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
	"github.com/kodmain/kitsune/src/services/gateway/server"
)

func main() {
	daemon.Start(server.Handler)
}
