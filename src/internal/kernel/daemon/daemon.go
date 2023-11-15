package daemon

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/multithread"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
)

const Name = "daemon"

type Handler struct {
	Name string
	Call func() error
}

type DaemonHandler struct {
	PIDHandler *PIDHandler
	sigs       chan os.Signal
	done       chan bool
}

func New() *DaemonHandler {
	return &DaemonHandler{
		PIDHandler: NewPIDHandler(config.BUILD_APP_NAME, config.PATH_RUN),
		sigs:       make(chan os.Signal, 1),
		done:       make(chan bool, 1),
	}
}

func (d *DaemonHandler) Start(handlers ...*Handler) {
	flag.Parse()
	if multithread.IsMaster() {
		pid := os.Getpid()
		if _, err := d.PIDHandler.GetPID(); err != nil {
			handleErrorAndExit(err)
		}

		if err := d.PIDHandler.SetPID(pid); err != nil {
			handleErrorAndExit(err)
		}
	}

	signal.Notify(d.sigs, syscall.SIGINT, syscall.SIGTERM)

	go d.handleSignal()

	for _, handler := range handlers {
		go d.processHandler(handler)
	}

	<-d.done
}

func (d *DaemonHandler) Stop() {
	d.sigs <- syscall.SIGTERM
}

func handleErrorAndExit(err error) {
	if logger.Fatal(err) {
		os.Exit(1)
	}
}

func (d *DaemonHandler) handleSignal() {
	<-d.sigs
	d.PIDHandler.ClearPID()
	d.done <- true
}

func (d *DaemonHandler) processHandler(handler *Handler) {
	count := 0
	startTime := time.Now()

	for {
		logger.Info(config.BUILD_APP_NAME + " " + handler.Name + " start")
		if err := handler.Call(); err != nil {
			logger.Warn(config.BUILD_APP_NAME+" "+handler.Name+" fail", err)
			if count >= 2 && d.shouldExit(count, startTime) {
				logger.Error(fmt.Errorf(config.BUILD_APP_NAME + " " + handler.Name + " won't start"))
				d.done <- true
				break
			}
			count++
		} else {
			break
		}
	}
}

func (d *DaemonHandler) shouldExit(count int, startTime time.Time) bool {
	elapsedTime := time.Since(startTime)
	return elapsedTime < time.Minute
}
