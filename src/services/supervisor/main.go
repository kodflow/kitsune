package main

import (
	"github.com/kodflow/kitsune/src/internal/kernel/daemon"
	"github.com/kodflow/kitsune/src/services/supervisor/process"
)

func main() {
	daemon.New().Start(process.Handler)
}
