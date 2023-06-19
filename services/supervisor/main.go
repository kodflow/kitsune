package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kodmain/KitsuneFramework/internal/env"
	"github.com/kodmain/KitsuneFramework/internal/kernel/daemon"
	"github.com/kodmain/KitsuneFramework/services/supervisor/process"
)

func main() {

	files, err := os.ReadDir(filepath.Join(env.PATH_LIB, env.PROJECT_NAME))
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du répertoire : %v\n", err)
		return
	}

	processes := make([]*process.Process, 0)
	for _, file := range files {
		if file.Name() != env.BUILD_APP_NAME {
			process := &process.Process{
				Name:    file.Name(),
				Command: filepath.Join(env.PATH_LIB, env.PROJECT_NAME, file.Name()),
				Restart: true,
			}
			processes = append(processes, process)
		}
	}

	for _, process := range processes {
		if err := process.Kill(); err == nil {
			process.Start()
		}
	}

	daemon.Start()
}
