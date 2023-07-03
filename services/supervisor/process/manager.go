package process

import (
	"fmt"
	"sync"

	"github.com/kodmain/KitsuneFramework/internal/kernel/daemon"
)

type Manager struct {
	processes map[string]*Process
	mutex     sync.Mutex
}

func NewProcessManager() *Manager {
	mngr := &Manager{
		processes: make(map[string]*Process),
		mutex:     sync.Mutex{},
	}

	return mngr
}

func (m *Manager) CreateProcess(name string, command string, args ...string) (*Process, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if process, _ := daemon.GetPID(name); process != nil {
		if err := daemon.ClearProcess(process, name); err != nil {
			return nil, err
		}
	}

	_, exists := m.processes[name]
	if exists {
		return nil, fmt.Errorf("process with name %s already exists", name)
	}

	proc := &Process{
		Name:    name,
		command: command,
		args:    args,
	}

	err := proc.Start()

	if err != nil {
		return nil, fmt.Errorf("failed to start process: %v", err)
	}

	m.processes[name] = proc
	return proc, nil
}

func (m *Manager) DeleteProcess(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	proc, exists := m.processes[name]
	if !exists {
		return fmt.Errorf("no process found with name %s", name)
	}

	err := proc.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	delete(m.processes, name)
	return nil
}

func (m *Manager) GetProcess(name string) (*Process, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	proc, exists := m.processes[name]
	return proc, exists
}

func (m *Manager) ListProcesses() map[string]*Process {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	copy := make(map[string]*Process, len(m.processes))
	for name, proc := range m.processes {
		copy[name] = proc
	}
	return copy
}
