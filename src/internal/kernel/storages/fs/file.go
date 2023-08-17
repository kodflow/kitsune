package fs

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CreateFile(filePath string, options ...*Options) (*os.File, error) {
	var opts *Options
	var err error

	if len(options) > 0 {
		opts = options[0]
	} else {
		opts, err = defaultOptions()
		if err != nil {
			return nil, err
		}
	}

	opts.fromFile = true

	err = CreateDirectory(filepath.Dir(filePath), opts)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	perms(filePath, opts)

	return file, nil
}

func OpenFile(filePath string, options ...*Options) (*os.File, error) {
	var opts *Options
	var err error

	if len(options) > 0 {
		opts = options[0]
	} else {
		opts, err = defaultOptions()
		if err != nil {
			return nil, err
		}
	}

	opts.fromFile = true

	err = CreateDirectory(filepath.Dir(filePath), opts)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, opts.Perms)
	if err != nil {
		return nil, err
	}

	perms(filePath, opts)

	return file, nil
}

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func SHA1File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
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

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func WriteFile(filePath string, content string) error {
	err := CreateDirectory(filepath.Dir(filePath))
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
