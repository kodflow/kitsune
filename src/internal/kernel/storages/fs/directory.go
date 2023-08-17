package fs

import (
	"os"
)

// Crée un dossier avec le chemin spécifié (récursivement si nécessaire)
func CreateDirectory(dirPath string, options ...*Options) error {
	var err error = nil

	if ExistsDirectory(dirPath) {
		return err
	}

	var opts *Options = nil

	if len(options) > 0 {
		opts = options[0]
		if opts.fromFile {
			opts = opts.Fork()
		}
	} else {
		opts, err = defaultOptions()
		if err != nil {
			return err
		}
	}

	opts.AddExecutable()

	err = os.MkdirAll(dirPath, opts.Perms)
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
