package daemon_test

import (
	"testing"
	"time"

	"github.com/kodflow/kitsune/src/internal/kernel/daemon"
)

func TestDaemonHandler_StartStop(t *testing.T) {
	handler := daemon.New()
	testHandler := &daemon.Handler{
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
