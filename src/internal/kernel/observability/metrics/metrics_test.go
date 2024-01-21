package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	var m *Metrics
	t.Run("New", func(t *testing.T) {
		m = New()
		assert.NotNil(t, m)
	})

	t.Run("GetDuration", func(t *testing.T) {
		// Test case 1: Duration metric exists
		duration := m.GetDuration("metric1")
		assert.NotNil(t, duration)
		assert.Equal(t, duration, m.GetDuration("metric1"))

		// Test case 2: Duration metric does not exist
		duration2 := m.GetDuration("metric2")
		assert.NotNil(t, duration2)
		assert.Equal(t, duration2, m.GetDuration("metric2"))
		assert.NotEqual(t, duration2, duration)
	})

	t.Run("GetCounter", func(t *testing.T) {
		// Test case 1: Counter metric exists
		counter := m.GetCounter("metric1")
		assert.NotNil(t, counter)
		assert.Equal(t, counter, m.GetCounter("metric1"))
		assert.Equal(t, counter.Value(), uint64(0))

		// Test case 2: Counter metric does not exist
		counter2 := m.GetCounter("metric2")
		counter2.Increment()
		assert.NotNil(t, counter2)
		assert.Equal(t, counter2, m.GetCounter("metric2"))
		assert.Equal(t, counter2.Value(), uint64(1))
		assert.NotEqual(t, counter2, counter)
	})

	t.Run("GetAverage", func(t *testing.T) {
		// Test case 1: Average metric exists
		average := m.GetAverage("metric1", time.Second)
		assert.NotNil(t, average)
		assert.Equal(t, average, m.GetAverage("metric1", time.Second))

		// Test case 2: Average metric does not exist
		average2 := m.GetAverage("metric2", time.Minute)
		assert.NotNil(t, average2)
		assert.Equal(t, average2, m.GetAverage("metric2", time.Minute))
		assert.NotEqual(t, average2, average)
	})
}
