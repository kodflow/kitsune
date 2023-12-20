package daemon

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaemonHandlerStartStop(t *testing.T) {
	handler := New()
	testHandler := &Handler{
		Name: "test",
		Call: func() error {
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

func TestDaemonHandlerStartStopError(t *testing.T) {
	handler := New()
	testHandler := &Handler{
		Name: "test",
		Call: func() error {
			return fmt.Errorf("test error")
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
			shouldExit: false,
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
		{
			name:       "GreaterThanMinute",
			count:      3,
			startTime:  startTime.Add(-time.Minute),
			shouldExit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			got := d.shouldExit(tt.count, tt.startTime)
			assert.Equal(t, tt.shouldExit, got)
		})
	}
}

func TestDaemonHandlerProcessHandlerFail(t *testing.T) {
	handler := &Handler{
		Name: "test",
		Call: func() error {
			return errors.New("test error")
		},
	}

	d := New()
	d.processHandler(handler)
	<-d.done
}
