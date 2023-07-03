package process

import (
	"os"
	"path/filepath"

	"github.com/kodmain/KitsuneFramework/internal/env"
	"github.com/kodmain/KitsuneFramework/internal/kernel/daemon"
)

var Handler *daemon.Handler = &daemon.Handler{
	Name: "process manager",
	Call: func() error {
		files, err := os.ReadDir(filepath.Join(env.PATH_SERVICES))
		if err != nil {
			return err
		}

		pm := NewProcessManager()
		for _, file := range files {
			if file.Name() != env.BUILD_APP_NAME {
				pm.CreateProcess(
					file.Name(),
					filepath.Join(env.PATH_SERVICES, file.Name()),
				)
			}
		}

		return nil
	},
}
