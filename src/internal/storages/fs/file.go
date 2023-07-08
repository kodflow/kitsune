package fs

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
)

// Crée un fichier avec le chemin spécifié
func CreateFile(filePath string, options ...*CreateOption) (*os.File, error) {
	var opts *CreateOption = nil

	if len(options) > 0 {
		opts = options[0]
	} else {
		u, _ := user.Current()
		g, _ := user.LookupGroupId(u.Gid)
		opts = &CreateOption{
			User:  u,
			Group: g,
			Perms: 0644,
		}
	}

	// Vérifie si le dossier parent existe, sinon le crée
	err := CreateDirectory(filepath.Dir(filePath), opts)
	if err != nil {
		return nil, err
	}

	// Crée le fichier
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	perms(filePath, opts)

	return file, nil
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
	err := CreateDirectory(filepath.Dir(filePath))
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
