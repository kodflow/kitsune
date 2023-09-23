package daemon_test

import (
	"testing"
	"time"

	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

func TestStartStop(t *testing.T) {
	handler := &daemon.Handler{
		Name: "TestHandler",
		Call: func() error {
			time.Sleep(2 * time.Second) // Simuler un travail qui prend du temps
			return nil
		},
	}

	// Lancez les handlers en arrière-plan
	go daemon.Start(handler)

	// Donnez un peu de temps pour démarrer le handler
	time.Sleep(1 * time.Second)

	// Stoppez les handlers
	daemon.Stop()

	// Vérifiez si le fichier PID a bien été supprimé
	if _, err := daemon.GetPID("TestHandler"); err == nil {
		t.Fatal("Expected error due to missing PID, got nil")
	}
}
