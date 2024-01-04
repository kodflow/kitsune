package probe

import (
	"sync/atomic"
)

// Counter represents a thread-safe counter.
// This structure provides a simple counter that can be safely incremented or decremented
// in a concurrent environment using atomic operations.
type Counter struct {
	current *uint64
}

// NewCounter initializes and returns a new instance of Counter.
// This function creates a Counter instance with the initial count set to zero.
//
// Returns:
// - *Counter: A pointer to the newly created Counter instance.
func NewCounter() *Counter {
	initialValue := uint64(0)
	return &Counter{
		current: &initialValue,
	}
}

// Increment safely increments the counter by one.
// This method uses atomic operations to ensure that the increment operation is thread-safe,
// preventing race conditions in a concurrent environment.
func (c *Counter) Increment() {
	atomic.AddUint64(c.current, 1)
}

// Decrement safely decrements the counter by one if it is greater than zero.
// This method uses atomic operations and a loop with a compare-and-swap operation to ensure
// that the decrement is thread-safe. The counter will not decrement below zero.
func (c *Counter) Decrement() {
	for {
		current := atomic.LoadUint64(c.current)
		if current == 0 {
			return // Counter cannot be negative
		}
		if atomic.CompareAndSwapUint64(c.current, current, current-1) {
			return
		}
	}
}

// Value returns the current value of the counter.
//
// This method uses atomic operations to ensure that the read operation is thread-safe.
// It returns the current value of the counter without modifying it.
//
// Returns:
// - uint64: The current value of the counter.
func (c *Counter) Value() uint64 {
	return atomic.LoadUint64(c.current)
}
