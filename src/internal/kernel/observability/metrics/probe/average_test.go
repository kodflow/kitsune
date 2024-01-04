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
		average = NewAverage(100 * time.Millisecond)
		assert.NotNil(t, average)
	})

	t.Run("Hit", func(t *testing.T) {
		average.Hit()
		assert.Equal(t, uint64(1), *average.hitCount)
		average.Hit()
		assert.Equal(t, uint64(2), *average.hitCount)
		average.Hit()
		assert.Equal(t, uint64(3), *average.hitCount)
	})

	t.Run("GetAverage", func(t *testing.T) {
		avg := average.GetAverage()
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, float64(3), avg)

		t.Run("Benchmark", func(t *testing.T) {
			for b := 0; b < 5; b++ {
				for i := 0; i < 1000; i++ {
					average.Hit()
				}
				time.Sleep(100 * time.Millisecond)
				v := average.GetAverage()
				assert.Greater(t, v, float64(900))
				assert.LessOrEqual(t, v, float64(1100))
			}
		})
	})
}
