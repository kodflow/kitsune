package main

import (
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
	"github.com/kodmain/kitsune/src/services/supervisor/process"
)

func main() {
	// update
	daemon.Start(process.Handler)
}
