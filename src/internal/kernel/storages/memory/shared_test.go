package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharedConcurrentOperations(t *testing.T) {
	shared := NewSharedMemory()
	const Concurrent = 100000
	var wg sync.WaitGroup

	t.Run("ConcurrentWrite", func(t *testing.T) {
		wg.Add(Concurrent)
		for i := 0; i < Concurrent; i++ {
			go func(i int) {
				defer wg.Done()
				shared.Store(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
			}(i)
		}
		wg.Wait()
	})

	t.Run("ConcurrentRead", func(t *testing.T) {
		wg.Add(Concurrent)
		for i := 0; i < Concurrent; i++ {
			go func(i int) {
				defer wg.Done()
				value, exists := shared.Read(fmt.Sprintf("key%d", i))
				assert.True(t, exists)
				assert.Equal(t, fmt.Sprintf("value%d", i), value)
			}(i)
		}
		wg.Wait()
	})

	t.Run("ConcurrentDelete", func(t *testing.T) {
		wg.Add(Concurrent)
		for i := 0; i < Concurrent; i++ {
			go func(i int) {
				defer wg.Done()
				shared.Delete(fmt.Sprintf("key%d", i))
				_, exists := shared.Read(fmt.Sprintf("key%d", i))
				assert.False(t, exists)
			}(i)
		}
		wg.Wait()
	})

	t.Run("ConcurrentExists", func(t *testing.T) {
		wg.Add(Concurrent)
		for i := 0; i < Concurrent; i++ {
			go func(i int) {
				defer wg.Done()
				exists := shared.Exists(fmt.Sprintf("key%d", i))
				assert.False(t, exists) // Les clés doivent être supprimées à ce stade
			}(i)
		}
		wg.Wait()
	})
}
