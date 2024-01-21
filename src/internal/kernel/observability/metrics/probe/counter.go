package probe

import (
	"sync/atomic"
)

// Counter represents a thread-safe counter.
// This structure provides a simple counter that can be safely incremented or decremented
// in a concurrent environment using atomic operations.
type Counter struct {
	value *uint64
}

// NewCounter initializes and returns a new instance of Counter.
// This function creates a Counter instance with the initial count set to zero.
//
// Returns:
// - *Counter: A pointer to the newly created Counter instance.
func NewCounter() *Counter {
	initialValue := uint64(0)
	return &Counter{
		value: &initialValue,
	}
}

// Increment safely increments the counter by one.
// This method uses atomic operations to ensure that the increment operation is thread-safe,
// preventing race conditions in a concurrent environment.
func (c *Counter) Increment() uint64 {
	atomic.AddUint64(c.value, 1)
	return c.Value()
}

// Decrement safely decrements the counter by one.
// This method uses atomic operations to ensure that the decrement operation is thread-safe,
func (c *Counter) Decrement() uint64 {
	if atomic.LoadUint64(c.value) > 0 {
		atomic.AddUint64(c.value, ^uint64(0))
	}

	return c.Value()
}

// Value returns the current value of the counter.
//
// This method uses atomic operations to ensure that the read operation is thread-safe.
// It returns the current value of the counter without modifying it.
//
// Returns:
// - uint64: The current value of the counter.
func (c *Counter) Value() uint64 {
	return atomic.LoadUint64(c.value)
}

func (c *Counter) Reset() {
	atomic.StoreUint64(c.value, 0)
}
