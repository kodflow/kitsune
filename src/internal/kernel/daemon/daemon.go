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

// Constant representing the name of the daemon.
const Name = "daemon"

// Handler struct defines a structure for handling specific daemon tasks.
type Handler struct {
	Name string       // Name of the handler.
	Call func() error // Call is the function to execute the handler's task.
}

// DaemonHandler manages the lifecycle and signal handling of the daemon.
type DaemonHandler struct {
	PIDHandler *PIDHandler    // PIDHandler is used to manage the PID file of the daemon.
	sigs       chan os.Signal // sigs is a channel for receiving OS signals.
	done       chan bool      // done is a channel to signal the completion of the daemon's execution.
}

// New creates a new instance of DaemonHandler.
// It initializes the PIDHandler and channels for signal and completion handling.
func New() *DaemonHandler {
	return &DaemonHandler{
		PIDHandler: NewPIDHandler(config.BUILD_APP_NAME, config.PATH_RUN),
		sigs:       make(chan os.Signal, 1),
		done:       make(chan bool, 1),
	}
}

// Start begins the daemon execution.
// It starts the specified handlers and listens for system signals to handle graceful shutdown.
// Parameters:
// - handlers: ...*Handler A variadic list of handlers to be executed by the daemon.
func (d *DaemonHandler) Start(handlers ...*Handler) {
	pid := os.Getpid()
	if _, err := d.PIDHandler.GetPID(); err != nil {
		handleErrorAndExit(err)
	}

	if err := d.PIDHandler.SetPID(pid); err != nil {
		handleErrorAndExit(err)
	}

	signal.Notify(d.sigs, syscall.SIGINT, syscall.SIGTERM)

	go d.handleSignal()

	for _, handler := range handlers {
		go d.processHandler(handler)
	}

	<-d.done
}

// Stop gracefully stops the daemon.
// It sends a SIGTERM signal to the daemon.
func (d *DaemonHandler) Stop() {
	d.sigs <- syscall.SIGTERM
}

// handleErrorAndExit logs the error and exits the program.
// Parameters:
// - err: error The error to log before exiting.
func handleErrorAndExit(err error) {
	if logger.Fatal(err) {
		os.Exit(1)
	}
}

// handleSignal waits for an OS signal and initiates the daemon shutdown process.
func (d *DaemonHandler) handleSignal() {
	<-d.sigs
	d.PIDHandler.ClearPID()
	d.done <- true
}

// processHandler executes a given handler and manages its lifecycle.
// It attempts to restart the handler on failure and stops the daemon if it fails repeatedly in a short time.
// Parameters:
// - handler: *Handler The handler to be processed.
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

// shouldExit determines if the daemon should exit based on the number of failures and time elapsed.
// Returns true if the handler has failed multiple times in a short duration.
// Parameters:
// - count: int The number of times the handler has failed.
// - startTime: time.Time The start time of the handler execution.
func (d *DaemonHandler) shouldExit(count int, startTime time.Time) bool {
	elapsedTime := time.Since(startTime)
	return elapsedTime < time.Minute
}
