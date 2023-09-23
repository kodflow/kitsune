package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

type Process struct {
	Name       string
	args       []string
	command    string
	IsRunning  bool
	Proc       *exec.Cmd
	pidHandler *daemon.PIDHandler
}

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
