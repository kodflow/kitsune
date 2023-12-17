package memory_test

import (
	"sync"
	"testing"

	"github.com/kodflow/kitsune/src/internal/kernel/storages/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRead(t *testing.T) {
	// Test case 1: Key exists in memory
	m := memory.New()
	m.Store("key1", "value1")

	value, exists := m.Read("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	// Test case 2: Key does not exist in memory
	value, exists = m.Read("key2")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestMemoryReadConcurrent(t *testing.T) {
	m := memory.New()
	m.Store("key1", "value1")

	var wg sync.WaitGroup
	numReaders := 100

	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			value, exists := m.Read("key1")
			assert.True(t, exists)
			assert.Equal(t, "value1", value)
		}()
	}

	wg.Wait()
}
