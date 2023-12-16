package fs

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"

	"github.com/kodmain/kitsune/src/config"
)

// Options represents the options for file system storage.
// This struct holds the configuration for file or directory access and permissions.
//
// Attributes:
// - User: *user.User The user that owns the file or directory.
// - Perms: fs.FileMode The permissions for the file or directory.
type Options struct {
	User  *user.User  // The user that owns the file.
	Perms fs.FileMode // The permissions for the file.
}

// defaultFileOptions returns the default options for the storage filesystem.
// It creates and returns a new Options struct with the default values for file storage.
// The default user is obtained from config.USER and the default permissions are set to 0644.
//
// Returns:
// - *Options: The default file Options.
// - error: nil since there is no error generation in this function.
func defaultFileOptions() (*Options, error) {
	return &Options{
		User:  config.USER,
		Perms: 0644,
	}, nil
}

// defaultDirectoryOptions returns the default options for a directory.
// It creates and returns a new Options struct with the default values for directory storage.
// The default user is config.USER and the default permissions are set to 0755.
//
// Returns:
// - *Options: The default directory Options.
// - error: nil since there is no error generation in this function.
func defaultDirectoryOptions() (*Options, error) {
	return &Options{
		User:  config.USER,
		Perms: 0755,
	}, nil
}

// resolveFileOptions resolves the file options.
// It checks if specific options are provided; if not, it defaults to the standard file options.
// This function is used to determine the final settings for file operations.
//
// Parameters:
// - options: []*Options Optional parameter for specifying custom file options.
//
// Returns:
// - *Options: The resolved file Options.
// - error: An error if resolving the options fails.
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
// It checks if specific options are provided; if not, it defaults to the standard directory options.
// This function is used to determine the final settings for directory operations.
//
// Parameters:
// - options: []*Options Optional parameter for specifying custom directory options.
//
// Returns:
// - *Options: The resolved directory Options.
// - error: An error if resolving the options fails.
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

// AddPerms adds the specified permissions to the current permissions of Options.
//
// Parameters:
// - perms: fs.FileMode The permissions to add to the current set.
func (co *Options) AddPerms(perms fs.FileMode) {
	co.Perms |= perms
}

// RemovePerms removes the specified permissions from the current permissions of Options.
//
// Parameters:
// - perms: fs.FileMode The permissions to remove from the current set.
func (co *Options) RemovePerms(perms fs.FileMode) {
	co.Perms &^= perms
}

// perms sets the permissions and ownership for a given path according to the provided Options.
// This function changes the owner and permissions of the file or directory at the given path.
//
// Parameters:
// - path: string The file system path for which to set permissions and ownership.
// - options: *Options The options specifying the desired permissions and owner.
//
// Returns:
// - error: An error if setting permissions or ownership fails.
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
