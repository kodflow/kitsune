package fs

import (
	"io/fs"
	"os"
)

// Crée un fichier avec le chemin spécifié
func CreateFile(filePath string) error {
	// Vérifie si le dossier parent existe, sinon le crée
	err := createParentDirectory(filePath)
	if err != nil {
		return err
	}

	// Crée le fichier
	_, err = os.Create(filePath)
	if err != nil {
		return err
	}

	return nil
}

// Supprime un fichier avec le chemin spécifié
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func ExistsFile(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func StatFile(filePath string) (fs.FileInfo, error) {
	return os.Stat(filePath)
}

// Lit le contenu d'un fichier avec le chemin spécifié
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// Écrit le contenu dans un fichier avec le chemin spécifié
func WriteFile(filePath string, content string) error {
	// Vérifie si le dossier parent existe, sinon le crée
	err := createParentDirectory(filePath)
	if err != nil {
		return err
	}

	// Ouvre le fichier en mode écriture
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrit le contenu dans le fichier
	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}
