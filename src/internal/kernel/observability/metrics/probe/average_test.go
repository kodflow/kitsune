package probe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAverage(t *testing.T) {
	var average *Average
	t.Run("NewAverage", func(t *testing.T) {
		assert.Nil(t, average)
		average = NewAverage(time.Minute * 2)
		average = NewAverage(100*time.Millisecond, time.Second*1, time.Second*3, time.Second*2)
		assert.NotNil(t, average)
		assert.Equal(t, uint64(30), average.capacity)
	})

	t.Run("Hit", func(t *testing.T) {
		for i := 1; i <= int(average.capacity+1); i++ {
			average.Hit()
			assert.Equal(t, uint64(i), average.counter.Value())
		}
	})

	t.Run("Reset", func(t *testing.T) {
		assert.Equal(t, uint64(31), average.counter.Value())
		average.Reset()
		assert.Equal(t, uint64(0), average.counter.Value())
	})

	t.Run("Value", func(t *testing.T) {
		for i := 0; i < 30; i++ {
			for b := 0; b < 10*i; b++ {
				average.Hit()
			}
			time.Sleep(100 * time.Millisecond)
		}
		avg := average.Value()
		assert.IsType(t, []float64{}, avg)
		assert.Equal(t, 3, len(avg))
		assert.Equal(t, 245.0, avg[0])
		assert.Equal(t, 145.0, avg[1])
		assert.Equal(t, 195.0, avg[2])
	})
}
