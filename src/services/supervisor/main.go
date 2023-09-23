package main

import (
	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
	"github.com/kodmain/kitsune/src/services/supervisor/process"
)

func main() {
	daemon.New(
		config.BUILD_APP_NAME,
		config.PATH_RUN,
	).Start(process.Handler)
}
