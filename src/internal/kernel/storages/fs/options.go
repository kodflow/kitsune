package fs

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"

	"github.com/kodmain/kitsune/src/config"
)

type Options struct {
	User  *user.User
	Perms fs.FileMode
}

func defaultOptions() (*Options, error) {
	return &Options{
		User:  config.USER,
		Perms: 0644,
	}, nil
}

func resolveOptions(options ...*Options) (*Options, error) {
	if len(options) > 0 && options[0] != nil {
		return options[0], nil
	}
	return defaultOptions()
}

func (co *Options) AddPerms(perms fs.FileMode) {
	co.Perms |= perms
}

func (co *Options) RemovePerms(perms fs.FileMode) {
	co.Perms &^= perms
}

func perms(path string, options *Options) error {
	uid, err := strconv.Atoi(options.User.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(options.User.Gid)
	if err != nil {
		return err
	}

	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}

	err = os.Chmod(path, options.Perms)
	if err != nil {
		return err
	}

	return nil
}
