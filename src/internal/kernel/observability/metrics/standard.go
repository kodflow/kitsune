package metrics

import (
	"sync"
	"time"

	"github.com/kodflow/kitsune/src/internal/kernel/observability/metrics/probe"
)

var (
	instance *Metrics = nil
	once     sync.Once
)

// standard returns an instance of the standard logger.
// If no instance exists, it is created with default parameters using New function.
// This function implements the singleton pattern to ensure only one instance of Logger exists.
//
// Returns:
// - *Logger: The singleton instance of the Logger.
func standard() *Metrics {
	once.Do(func() {
		instance = New()
	})
	return instance
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
func GetDuration(name string) *probe.Duration {
	return standard().GetDuration(name)
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
func GetCounter(name string) *probe.Counter {
	return standard().GetCounter(name)
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
func GetAverage(name string, duration time.Duration) *probe.Average {
	return standard().GetAverage(name, duration)
}
