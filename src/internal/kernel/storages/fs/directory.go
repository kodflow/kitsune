package fs

import "os"

func CreateDirectory(dirPath string, options ...*Options) error {
	if ExistsDirectory(dirPath) {
		return nil
	}

	opts, err := resolveOptions(options...)
	if err != nil {
		return err
	}

	opts.AddPerms(0111)

	err = os.MkdirAll(dirPath, opts.Perms)
	if err != nil {
		return err
	}

	return perms(dirPath, opts)
}

func ExistsDirectory(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func DeleteDirectory(dirPath string) error {
	return os.RemoveAll(dirPath)
}
