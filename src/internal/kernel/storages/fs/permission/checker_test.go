package permission

import (
	"os"
	"testing"
)

func TestCheck_Success(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Définir les permissions pour le dossier
	err = os.Chmod(dir, 0700)
	if err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}

	// Vérifier les permissions en utilisant la fonction Check
	err = Check(dir, 0700)
	if err != nil {
		t.Errorf("failed to check permissions: %v", err)
	}
}

func TestCheck_Fail(t *testing.T) {
	// Utilisation d'un chemin non existant avec des permissions aléatoires pour échouer le test
	if err := Check("/fake/path", 0700); err == nil {
		t.Errorf("expected to fail permission check, but passed")
	} else {
		t.Log("expected failure, error:", err)
	}
}
func TestCheck_InvalidPermissions(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-permission")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Définir les permissions pour le dossier
	err = os.Chmod(dir, 0755) // permissions de lecture/écriture pour l'utilisateur, lecture/exécution pour les autres
	if err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}

	// Essayer de vérifier le répertoire avec des permissions qui ne sont pas celles définies
	err = Check(dir, 0700) // 0700 = rwx pour l'utilisateur uniquement
	if err == nil {
		currentPerms, _ := os.Stat(dir)
		t.Errorf("expected to fail permission check for invalid permissions, but passed. Current permissions: %v", currentPerms.Mode().Perm())
	} else {
		t.Log("expected failure for invalid permissions, error:", err)
	}
}
