package fs

import (
	"os"
	"path/filepath"
)

// Crée un dossier avec le chemin spécifié (récursivement si nécessaire)
func CreateDirectory(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func ExistsDirectory(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Supprime un dossier avec le chemin spécifié (y compris les sous-dossiers et les fichiers)
func DeleteDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}

	return nil
}

// Fonction privée pour créer le dossier parent d'un fichier si nécessaire
func createParentDirectory(filePath string) error {
	parentDir := filepath.Dir(filePath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		err := os.MkdirAll(parentDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
