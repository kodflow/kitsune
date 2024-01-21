package probe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDuration(t *testing.T) {
	var duration *Duration

	t.Run("NewDuration", func(t *testing.T) {
		assert.Nil(t, duration)
		duration = NewDuration()
		assert.NotNil(t, duration)

		currentTime := time.Now()
		startTime := time.Unix(0, int64(*duration.start))
		assert.True(t, startTime.Before(currentTime) || startTime.Equal(currentTime))
	})

	t.Run("start", func(t *testing.T) {
		startTime := time.Unix(0, int64(*duration.start))
		duration.Start()

		currentTime := time.Now()
		updatedStartTime := time.Unix(0, int64(*duration.start))
		assert.True(t, updatedStartTime.Before(currentTime) || updatedStartTime.Equal(currentTime))
		assert.NotEqual(t, startTime, updatedStartTime)
	})

	t.Run("elapsed", func(t *testing.T) {
		elapsedTime := duration.Elapsed()
		assert.True(t, elapsedTime > 0)
	})

}
