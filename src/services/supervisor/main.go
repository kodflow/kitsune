package main

import (
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
	"github.com/kodmain/kitsune/src/services/supervisor/process"
)

func main() {
	daemon.Start(process.Handler)
}
