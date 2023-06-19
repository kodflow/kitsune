package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kodmain/KitsuneFramework/internal/env"
	"github.com/kodmain/KitsuneFramework/internal/kernel/daemon"
)

type Process struct {
	Name    string
	Command string
	Args    []string
	Restart bool
	Running bool
	Proc    *exec.Cmd
}

func (p *Process) Kill() error {
	pidFile := filepath.Join(env.PATH_PID, p.Name+".pid")
	if process, _ := daemon.GetPID(pidFile); process != nil {
		err := process.Kill()
		if err != nil {
			return fmt.Errorf("failed to kill process: %v", err)
		}

		err = os.Remove(pidFile)
		if err != nil {
			return fmt.Errorf("failed to remove PID file: %v", err)
		}
	}

	return nil
}

func (p *Process) Start() error {
	if p.Running {
		return fmt.Errorf("process %s is already running", p.Name)
	}

	cmd := exec.Command(p.Command, p.Args...)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start process %s: %v", p.Name, err)
	}

	p.Running = true
	p.Proc = cmd

	go func() {
		_ = cmd.Wait()
		p.Running = false

		if p.Restart {
			fmt.Printf("process %s has stopped, restarting...\n", p.Name)
			p.Start()
		} else {
			fmt.Printf("process %s has stopped\n", p.Name)
		}
	}()

	return nil
}

func (p *Process) Stop() error {
	if !p.Running {
		return fmt.Errorf("process %s is not running", p.Name)
	}

	err := p.Proc.Process.Kill()
	if err != nil {
		return fmt.Errorf("failed to stop process %s: %v", p.Name, err)
	}

	p.Running = false
	return nil
}
