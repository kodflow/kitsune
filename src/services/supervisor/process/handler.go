package process

import (
	"os"
	"path/filepath"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

// Handler represents the process manager handler.
var Handler *daemon.Handler = &daemon.Handler{
	Name: "process manager",
	Call: func() error {
		files, err := os.ReadDir(filepath.Join(config.PATH_SERVICES))
		if err != nil {
			return err
		}

		pm := NewProcessManager()
		for _, file := range files {
			if file.Name() != config.BUILD_APP_NAME {
				pm.CreateProcess(
					file.Name(),
					filepath.Join(config.PATH_SERVICES, file.Name()),
				)
			}
		}

		return nil
	},
}
