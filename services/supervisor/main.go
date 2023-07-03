package main

import (
	"github.com/kodmain/kitsune/internal/kernel/daemon"
	"github.com/kodmain/kitsune/services/supervisor/process"
)

func main() {
	daemon.Start(process.Handler)
}
