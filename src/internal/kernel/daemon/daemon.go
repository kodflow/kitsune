package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/kodmain/kitsune/src/internal/observability/logger"
)

type Handler struct {
	Name string
	Call func() error
}

var sigs chan os.Signal = make(chan os.Signal, 1)
var done chan bool = make(chan bool, 1)

func Start(handlers ...*Handler) {
	fmt.Println(env.BUILD_APP_NAME, "start")

	if _, err := GetPID(env.BUILD_APP_NAME); err != nil {
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
	fmt.Println(env.BUILD_APP_NAME, "exit")
}

func Stop() {
	fmt.Println(env.BUILD_APP_NAME, "stop")
	sigs <- syscall.SIGTERM
}

func handleErrorAndExit(err error) {
	if logger.Default().Fatal(err) {
		os.Exit(1)
	}
}

func handleSignal() {
	<-sigs
	ClearPID(env.BUILD_APP_NAME)
	done <- true
}

func processHandler(handler *Handler) {
	count := 0
	startTime := time.Now()

	for {
		logger.Default().Info(env.BUILD_APP_NAME + " " + handler.Name + " start")
		if err := handler.Call(); err != nil {
			logger.Default().Warn(env.BUILD_APP_NAME + " " + handler.Name + " fail")
			if count >= 2 && shouldExit(count, startTime) {
				logger.Default().Error(fmt.Errorf(env.BUILD_APP_NAME + " " + handler.Name + " won't start"))
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
