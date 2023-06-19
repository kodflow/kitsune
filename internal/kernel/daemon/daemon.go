package daemon

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kodmain/KitsuneFramework/internal/kernel/exit"
	"github.com/kodmain/KitsuneFramework/internal/storages/fs"
)

func Start() {
	if _, err := GetPID(pidFilePath); err != nil {
		exit.App(exit.Error, err.Error())
	}

	if err := SetPID(pidFilePath); err != nil {
		exit.App(exit.Error, err.Error())
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fs.DeleteFile(pidFilePath)
		done <- true
	}()
	<-done
}
