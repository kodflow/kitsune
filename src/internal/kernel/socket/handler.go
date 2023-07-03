package socket

import (
	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

var Handler *daemon.Handler = &daemon.Handler{
	Name: "Socket Server",
	Call: func() error {
		return Server(env.PATH_RUN + env.BUILD_APP_NAME).Start()
	},
}
