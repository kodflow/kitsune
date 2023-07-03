package main

import (
	"github.com/kodmain/KitsuneFramework/internal/kernel/daemon"
	"github.com/kodmain/KitsuneFramework/services/supervisor/process"
)

func main() {
	daemon.Start(process.Handler)
}
