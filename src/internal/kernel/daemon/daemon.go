package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
)

type Handler struct {
	Name string
	Call func() error
}

var sigs chan os.Signal = make(chan os.Signal, 1)
var done chan bool = make(chan bool, 1)

func Start(handlers ...*Handler) {
	if _, err := GetPID(config.BUILD_APP_NAME); err != nil {
		handleErrorAndExit(err)
	}

	if err := SetPID(); err != nil {
		handleErrorAndExit(err)
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go handleSignal()

	for _, handler := range handlers {
		go processHandler(handler)
	}

	<-done
}

func Stop() {
	sigs <- syscall.SIGTERM
}

func handleErrorAndExit(err error) {
	if logger.Fatal(err) {
		os.Exit(1)
	}
}

func handleSignal() {
	<-sigs
	ClearPID(config.BUILD_APP_NAME)
	done <- true
}

func processHandler(handler *Handler) {
	count := 0
	startTime := time.Now()

	for {
		logger.Info(config.BUILD_APP_NAME + " " + handler.Name + " start")
		if err := handler.Call(); err != nil {
			logger.Warn(config.BUILD_APP_NAME + " " + handler.Name + " fail")
			if count >= 2 && shouldExit(count, startTime) {
				logger.Error(fmt.Errorf(config.BUILD_APP_NAME + " " + handler.Name + " won't start"))
				done <- true
				break
			}
			count++
		} else {
			break
		}
	}
}

func shouldExit(count int, startTime time.Time) bool {
	elapsedTime := time.Since(startTime)
	return elapsedTime < time.Minute
}
