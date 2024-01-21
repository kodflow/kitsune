package probe

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	var counter *Counter
	t.Run("NewCounter", func(t *testing.T) {
		assert.Nil(t, counter)
		counter = NewCounter()
		assert.NotNil(t, counter)
	})

	t.Run("Increment", func(t *testing.T) {
		counter.Increment()
		counter.Increment()
		assert.Equal(t, uint64(2), counter.Value(), "Expected counter value to be 1 after incrementing")
	})

	t.Run("Decrement", func(t *testing.T) {
		counter.Decrement()
		assert.Equal(t, uint64(1), counter.Value(), "Expected counter value to be 0 after decrementing")
	})

	t.Run("ConcurrentIncrement", func(t *testing.T) {
		counter.Reset()
		assert.Equal(t, uint64(0), counter.Value(), "Expected counter value to be 0 after reset")
		var wg sync.WaitGroup
		numRoutines := 100
		expectedValue := uint64(numRoutines)

		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go func() {
				counter.Increment()
				wg.Done()
			}()
		}

		wg.Wait()
		assert.Equal(t, expectedValue, counter.Value(), "Expected counter value to be %d after concurrent increment", expectedValue)
	})

}
