package process

import (
	"fmt"
	"sync"

	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/kernel/daemon"
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

// CreateProcess creates a new process with the given name, command, and arguments.
// It returns a pointer to the created Process and an error if any.
// If a process with the same name is already running or exists, an error is returned.
func (m *Manager) CreateProcess(name string, command string, args ...string) (*Process, error) {
	pidHandler := daemon.NewPIDHandler(name, config.PATH_RUN)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if pid, _ := pidHandler.GetPID(); pid != 0 {
		return nil, fmt.Errorf("process with name %s is already running", name)
	}

	_, exists := m.processes[name]
	if exists {
		return nil, fmt.Errorf("process with name %s already exists", name)
	}

	proc := &Process{
		Name:       name,
		command:    command,
		args:       args,
		pidHandler: pidHandler,
	}

	err := proc.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start process: %v", err)
	}

	m.processes[name] = proc
	return proc, nil
}

// DeleteProcess deletes a process with the given name from the manager.
// It kills the process and removes it from the list of managed processes.
// If no process is found with the given name, it returns an error.
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

// GetProcess récupère le processus avec le nom spécifié.
// Il renvoie le processus et un booléen indiquant si le processus existe.
func (m *Manager) GetProcess(name string) (*Process, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	proc, exists := m.processes[name]
	return proc, exists
}

// ListProcesses returns a map of all processes managed by the Manager.
// The key of the map is the name of the process, and the value is a pointer to the Process struct.
func (m *Manager) ListProcesses() map[string]*Process {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	copy := make(map[string]*Process, len(m.processes))
	for name, proc := range m.processes {
		copy[name] = proc
	}
	return copy
}
