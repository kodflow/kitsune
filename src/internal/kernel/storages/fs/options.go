package fs

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"

	"github.com/kodmain/kitsune/src/internal/env"
)

type Options struct {
	User     *user.User
	Perms    fs.FileMode
	fromFile bool
}

func defaultOptions() (*Options, error) {
	return &Options{
		User:  env.USER,
		Perms: 0644,
	}, nil
}

func (co *Options) Fork() *Options {
	return &Options{
		User:  co.User,
		Perms: co.Perms,
	}
}

func (co *Options) AddRead() {
	co.Perms |= 0444
}

func (co *Options) AddWrite() {
	co.Perms |= 0222
}

func (co *Options) AddExecutable() {
	co.Perms |= 0111
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
