package fs

import (
	"os"
	"os/user"
)

// Crée un dossier avec le chemin spécifié (récursivement si nécessaire)
func CreateDirectory(dirPath string, options ...*CreateOption) error {
	if ExistsDirectory(dirPath) {
		return nil
	}

	var opts *CreateOption = nil
	if len(options) > 0 {
		opts = options[0]
	} else {
		u, _ := user.Current()
		g, _ := user.LookupGroupId(u.Gid)
		opts = &CreateOption{
			User:  u,
			Group: g,
			Perms: 0755,
		}
	}

	err := os.MkdirAll(dirPath, opts.Perms)
	if err != nil {
		return err
	}

	return perms(dirPath, opts)
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
