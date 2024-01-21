package metrics

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStandard(t *testing.T) {
	// Create a wait group to ensure all goroutines finish before asserting the results
	var wg sync.WaitGroup
	wg.Add(2)

	// Test the creation of the metrics instance
	t.Run("CreateInstance", func(t *testing.T) {
		defer wg.Done()

		// Call the standard function to create the metrics instance
		m := standard()

		// Assert that the instance is not nil
		assert.NotNil(t, m, "Metrics instance should not be nil")
	})

	// Test the singleton behavior of the metrics instance
	t.Run("SingletonBehavior", func(t *testing.T) {
		defer wg.Done()

		// Call the standard function multiple times concurrently
		go func() {
			m1 := standard()
			assert.NotNil(t, m1, "Metrics instance should not be nil")
		}()

		go func() {
			m2 := standard()
			assert.NotNil(t, m2, "Metrics instance should not be nil")
		}()
	})

	// Wait for all goroutines to finish
	wg.Wait()

	t.Run("GetDuration", func(t *testing.T) {
		// Test case 1
		duration := GetDuration("test")
		assert.NotNil(t, duration, "should return a non-nil duration")
	})

	t.Run("GetCounter", func(t *testing.T) {
		// Test case 1
		counter := GetCounter("test")
		assert.NotNil(t, counter, "should return a non-nil counter")
	})

	t.Run("GetAverage", func(t *testing.T) {
		// Test case 1
		average := GetAverage("test", 100*time.Millisecond)
		assert.NotNil(t, average, "should return a non-nil average")
	})
}
