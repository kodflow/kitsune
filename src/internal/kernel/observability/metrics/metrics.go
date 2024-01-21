package metrics

import (
	"sync"
	"time"

	"github.com/kodflow/kitsune/src/internal/kernel/observability/metrics/probe"
)

// Metrics structure holds various types of metrics like counters, durations, and averages.
//
// Metrics provides a thread-safe way to create and manage different types of metric probes
// such as counters, durations, and averages. It uses sync.RWMutex to handle concurrent access.
type Metrics struct {
	counters  map[string]*probe.Counter
	durations map[string]*probe.Duration
	averages  map[string]*probe.Average
	mu        sync.RWMutex
}

// New creates and returns a new Metrics instance.
//
// New initializes a Metrics instance with empty maps for counters, durations, and averages.
//
// Returns:
// - *Metrics: Pointer to the newly created Metrics instance.
func New() *Metrics {
	return &Metrics{
		counters:  make(map[string]*probe.Counter),
		durations: make(map[string]*probe.Duration),
		averages:  make(map[string]*probe.Average),
	}
}

// GetDuration retrieves or creates a duration metric.
//
// GetDuration first attempts to find an existing duration metric by its name. If it doesn't exist,
// a new one is created and added to the map. It ensures thread safety using mutex locks.
//
// Parameters:
// - name: string - The name of the duration metric.
//
// Returns:
// - *probe.Duration: Pointer to the retrieved or newly created duration metric.
func (m *Metrics) GetDuration(name string) *probe.Duration {
	m.mu.RLock()
	if metric, ok := m.durations[name]; ok {
		m.mu.RUnlock()
		return metric
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// thread-safe check
	if metric, ok := m.durations[name]; ok {
		return metric
	}

	metric := probe.NewDuration()
	m.durations[name] = metric

	return metric
}

// GetCounter retrieves or creates a counter metric.
//
// Similar to GetDuration, GetCounter looks for an existing counter metric or creates a new one
// if it doesn't exist, ensuring thread safety with mutex locks.
//
// Parameters:
// - name: string - The name of the counter metric.
//
// Returns:
// - *probe.Counter: Pointer to the retrieved or newly created counter metric.
func (m *Metrics) GetCounter(name string) *probe.Counter {
	m.mu.RLock()
	if metric, ok := m.counters[name]; ok {
		m.mu.RUnlock()
		return metric
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// thread-safe check
	if metric, ok := m.counters[name]; ok {
		return metric
	}

	metric := probe.NewCounter()
	m.counters[name] = metric

	return metric
}

// GetAverage retrieves or creates an average metric.
//
// GetAverage operates like GetDuration and GetCounter, handling average metrics.
// It additionally requires a duration parameter for the new average metric.
//
// Parameters:
// - name: string - The name of the average metric.
// - duration: time.Duration - The duration for the average calculation.
//
// Returns:
// - *probe.Average: Pointer to the retrieved or newly created average metric.
func (m *Metrics) GetAverage(name string, duration time.Duration) *probe.Average {
	m.mu.RLock()
	if metric, ok := m.averages[name]; ok {
		m.mu.RUnlock()
		return metric
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// thread-safe check
	if metric, ok := m.averages[name]; ok {
		return metric
	}

	metric := probe.NewAverage(duration)
	m.averages[name] = metric

	return metric
}
