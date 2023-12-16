package fs

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"

	"github.com/kodmain/kitsune/src/config"
)

// Options represents the options for file system storage.
type Options struct {
	User  *user.User  // The user that owns the file.
	Perms fs.FileMode // The permissions for the file.
}

// defaultOptions returns the default options for the storage filesystem.
// It creates and returns a new Options struct with the default values.
// The default values include the user and permissions.
// The user is obtained from the config.USER constant.
// The permissions are set to 0644.
func defaultFileOptions() (*Options, error) {
	return &Options{
		User:  config.USER,
		Perms: 0644,
	}, nil
}

// defaultDirectoryOptions returns the default options for a directory.
// It creates and returns a new Options struct with the following values:
// - User: the value of config.USER
// - Perms: 0755
// It returns the newly created Options struct and nil error.
func defaultDirectoryOptions() (*Options, error) {
	return &Options{
		User:  config.USER,
		Perms: 0755,
	}, nil
}

// resolveFileOptions resolves the file options.
// It returns the provided file options if any, or the default options if no options are specified.
//
// Parameters:
// - options: The file options (optional)
//
// Returns:
// - *Options: The resolved file options
// - error: An error if there was a problem resolving the options
func resolveFileOptions(options ...*Options) (*Options, error) {
	if len(options) > 0 && options[0] != nil {
		if options[0].User == nil {
			options[0].User = config.USER
		}

		if options[0].Perms == 0 {
			options[0].Perms = 0644
		}

		return options[0], nil
	}

	return defaultFileOptions()
}

// resolveDirectoryOptions resolves the directory options.
// It returns the provided directory options if any, or the default options if no options are specified.
//
// Parameters:
// - options: The directory options (optional)
//
// Returns:
// - *Options: The resolved directory options
// - error: An error if there was a problem resolving the options
func resolveDirectoryOptions(options ...*Options) (*Options, error) {
	if len(options) > 0 && options[0] != nil {
		if options[0].User == nil {
			options[0].User = config.USER
		}

		if options[0].Perms == 0 {
			options[0].Perms = 0755
		}

		return options[0], nil
	}

	return defaultDirectoryOptions()
}

// AddPerms ajoute les droits d'accès spécifiés aux options de stockage.
// Les droits d'accès spécifiés sont ajoutés aux droits d'accès existants.
//
// Paramètres :
// - perms : Les droits d'accès à ajouter.
//
// Retour :
// Aucun.
func (co *Options) AddPerms(perms fs.FileMode) {
	co.Perms |= perms
}

// RemovePerms supprime les permissions spécifiées des options.
// Les permissions à supprimer sont spécifiées par le paramètre "perms".
// Cette méthode met à jour les permissions des options en utilisant l'opérateur
// bitwise AND NOT (^=) pour supprimer les permissions spécifiées.
//
// Paramètres:
// - perms: Les permissions à supprimer.
//
// Retour:
// Aucun.
func (co *Options) RemovePerms(perms fs.FileMode) {
	co.Perms &^= perms
}

// perms sets the permissions and ownership of the specified path.
// It takes the path and options as parameters and returns an error if any.
//
// Parameters:
// - path: The path to set permissions and ownership for
// - options: The options containing the desired permissions and ownership
//
// Returns:
// - error: An error if any of the operations fail
func perms(path string, options *Options) error {
	fmt.Println("perms", options)
	uid, err := strconv.Atoi(options.User.Uid)
	if err != nil {
		fmt.Println("a")
		return err
	}

	gid, err := strconv.Atoi(options.User.Gid)
	if err != nil {
		fmt.Println("b")
		return err
	}

	err = os.Chown(path, uid, gid)
	if err != nil {
		fmt.Println("c")
		return err
	}

	err = os.Chmod(path, options.Perms)
	if err != nil {
		fmt.Println("d")
		return err
	}

	fmt.Println("e")
	return nil
}
