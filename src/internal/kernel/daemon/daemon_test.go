package daemon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaemonHandler_StartStop(t *testing.T) {
	handler := New()
	testHandler := &Handler{
		Name: "test",
		Call: func() error {
			time.Sleep(10 * time.Second)
			return nil
		},
	}

	time.AfterFunc(0*time.Second, func() {
		time.AfterFunc(3*time.Second, func() {
			t.Error("Daemon did not stop in the expected time")
		})
		handler.Stop()
	})

	handler.Start(testHandler)
}
func TestDaemonHandlerShouldExit(t *testing.T) {
	startTime := time.Now()

	tests := []struct {
		name       string
		count      int
		startTime  time.Time
		shouldExit bool
	}{
		{
			name:       "NoFailures",
			count:      0,
			startTime:  startTime,
			shouldExit: true,
		},
		{
			name:       "LessThanMinute",
			count:      3,
			startTime:  startTime,
			shouldExit: true,
		},
		{
			name:       "GreaterThanMinute",
			count:      3,
			startTime:  startTime.Add(-time.Minute),
			shouldExit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DaemonHandler{}
			got := d.shouldExit(tt.count, tt.startTime)
			assert.Equal(t, tt.shouldExit, got)
		})
	}
}
