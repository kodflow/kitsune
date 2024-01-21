package probe

import (
	"sync/atomic"
	"time"
)

// Duration measures the elapsed time from a starting point.
// This structure allows for the recording of elapsed time in a thread-safe manner
// using atomic operations for the start time.
type Duration struct {
	start *uint64
}

// NewDuration initializes and returns a new instance of Duration.
// It sets the start time to the current time at the moment of creation.
//
// Returns:
// - *Duration: A pointer to the newly created Duration instance.
func NewDuration() *Duration {
	startVal := uint64(time.Now().UnixNano())

	return &Duration{
		start: &startVal,
	}
}

// Elapsed returns the elapsed time since start.
// The method calculates the time elapsed since the start time, which is stored as Unix nanoseconds,
// and returns the duration in time.Duration format.
//
// Returns:
// - time.Duration: The elapsed time since start as a time.Duration.
func (pd *Duration) Elapsed() time.Duration {
	start := atomic.LoadUint64(pd.start)
	now := uint64(time.Now().UnixNano())
	return time.Duration(now - start)
}

// Start sets or resets the start time of the duration measurement.
// This method updates the start time to the current time using atomic operations
// to ensure thread safety.
func (pd *Duration) Start() {
	atomic.StoreUint64(pd.start, uint64(time.Now().UnixNano()))
}
