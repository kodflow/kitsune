package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

// Process represents a process that can be started and managed.
type Process struct {
	Name       string             // Name is the name of the process.
	args       []string           // args are the arguments passed to the process.
	command    string             // command is the command used to start the process.
	IsRunning  bool               // IsRunning indicates whether the process is currently running.
	Proc       *exec.Cmd          // Proc is the underlying exec.Cmd instance representing the process.
	pidHandler *daemon.PIDHandler // pidHandler is used to handle the process ID.
}

// Kill terminates the process associated with the Process object.
// It first retrieves the process ID (PID) using the pidHandler.
// If the PID is not 0, it uses os.FindProcess to find the process.
// Then, it calls the Kill method on the process to terminate it.
// If an error occurs during the process termination or while clearing the PID file,
// it returns an error with a descriptive message.
// If the PID is 0, indicating that no process is associated with the Process object,
// it returns nil.
func (p *Process) Kill() error {
	if pid, _ := p.pidHandler.GetPID(); pid != 0 {
		process, err := os.FindProcess(pid)
		if err != nil {
			return err
		}

		err = process.Kill()
		if err != nil {
			return fmt.Errorf("failed to kill process: %v", err)
		}

		err = p.pidHandler.ClearPID()
		if err != nil {
			return fmt.Errorf("failed to remove PID file: %v", err)
		}
	}

	return nil
}

// Restart restarts the process.
// It stops the process, waits for 1 second, and then starts it again.
// If stopping the process or starting it again fails, an error is returned.
func (p *Process) Restart() error {
	err := p.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop process for restart: %v", err)
	}

	time.Sleep(1 * time.Second)

	err = p.Start()
	if err != nil {
		return fmt.Errorf("failed to restart process: %v", err)
	}

	return nil
}

// Start starts the process.
// It returns an error if the process is already running or if it fails to start.
// The function starts the process in a separate goroutine and updates the IsRunning flag accordingly.
// If the process stops, it will be automatically restarted.
func (p *Process) Start() error {
	if p.IsRunning {
		return fmt.Errorf("process %s is already Isrunning", p.Name)
	}

	cmd := exec.Command(p.command, p.args...)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start process %s: %v", p.Name, err)
	}

	p.IsRunning = true
	p.Proc = cmd

	go func() {
		_ = cmd.Wait()
		p.IsRunning = false
		fmt.Printf("process %s has stopped\n", p.Name)
		p.Start()
	}()

	return nil
}

// Stop stops the process.
// It returns an error if the process is not running or if there was an error while stopping the process.
func (p *Process) Stop() error {
	if !p.IsRunning {
		return fmt.Errorf("process %s is not Isrunning", p.Name)
	}

	err := p.Proc.Process.Signal(syscall.SIGTERM)

	if err != nil {
		return fmt.Errorf("failed to stop process %s: %v", p.Name, err)
	}

	p.IsRunning = false
	return nil
}
