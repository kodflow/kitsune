package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	m := NewMemory()
	assert.NotNil(t, m)

	t.Run("Instance", func(t *testing.T) {
		t.Run("Write", func(t *testing.T) {
			m.Store("key1", "value1")
		})

		t.Run("Read", func(t *testing.T) {
			value, exists := m.Read("key1")
			assert.True(t, exists)
			assert.Equal(t, "value1", value)
		})

		t.Run("Delete", func(t *testing.T) {
			m.Delete("key1")
			value, exists := m.Read("key1")
			assert.False(t, exists)
			assert.Nil(t, value)
		})
	})

	t.Run("Concurrent", func(t *testing.T) {
		const Concurrent = 100000
		var wg sync.WaitGroup
		t.Run("Write", func(t *testing.T) {
			wg.Add(Concurrent)
			for i := 0; i < Concurrent; i++ {
				go func(i int) {
					defer wg.Done()
					m.Store(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
				}(i)
			}
		})

		wg.Wait()
		t.Run("Read", func(t *testing.T) {
			wg.Add(Concurrent)
			for i := 0; i < Concurrent; i++ {
				go func(i int) {
					defer wg.Done()
					value, exists := m.Read(fmt.Sprintf("key%d", i))
					assert.True(t, exists)
					assert.Equal(t, fmt.Sprintf("value%d", i), value)
				}(i)
			}
		})

		wg.Wait()

		t.Run("Delete", func(t *testing.T) {
			wg.Add(Concurrent)
			for i := 0; i < Concurrent; i++ {
				go func(i int) {
					defer wg.Done()
					m.Delete(fmt.Sprintf("key%d", i))
					value, exists := m.Read(fmt.Sprintf("key%d", i))
					assert.False(t, exists)
					assert.Nil(t, value)
				}(i)
			}
		})

		wg.Wait()
	})
}
