package fs

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"
)

type CreateOption struct {
	User  *user.User
	Group *user.Group
	Perms fs.FileMode
}

func (co *CreateOption) AddRead() {
	co.Perms |= 0444 // Ajoute les permissions de lecture pour le propriétaire, le groupe et les autres
}

func (co *CreateOption) AddWrite() {
	co.Perms |= 0222 // Ajoute les permissions d'écriture pour le propriétaire, le groupe et les autres
}

func (co *CreateOption) AddExecutable() {
	co.Perms |= 0111 // Ajoute les permissions d'exécution pour le propriétaire, le groupe et les autres
}

func perms(path string, options *CreateOption) error {
	uid, err := strconv.Atoi(options.User.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(options.Group.Gid)
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
