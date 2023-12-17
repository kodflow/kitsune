package fs

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"

	"github.com/kodflow/kitsune/src/config"
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
func defaultFileOptions() *Options {
	return &Options{
		User:  config.USER,
		Perms: 0644,
	}
}

// defaultDirectoryOptions returns the default options for a directory.
// It creates and returns a new Options struct with the default values for directory storage.
// The default user is config.USER and the default permissions are set to 0755.
//
// Returns:
// - *Options: The default directory Options.
// - error: nil since there is no error generation in this function.
func defaultDirectoryOptions() *Options {
	return &Options{
		User:  config.USER,
		Perms: 0755,
	}
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
func resolveFileOptions(options ...*Options) *Options {
	if len(options) > 0 && options[0] != nil {
		if options[0].User == nil {
			options[0].User = config.USER
		}

		if options[0].Perms == 0 {
			options[0].Perms = 0644
		}

		return options[0]
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
func resolveDirectoryOptions(options ...*Options) *Options {
	if len(options) > 0 && options[0] != nil {
		if options[0].User == nil {
			options[0].User = config.USER
		}

		if options[0].Perms == 0 {
			options[0].Perms = 0755
		}

		return options[0]
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
	uid, err := parseUID(options)
	if err != nil {
		return err
	}

	gid, err := parseGID(options)
	if err != nil {
		return err
	}

	err = changeOwnership(path, uid, gid)
	if err != nil {
		return err
	}

	return changePermissions(path, options)
}

// parseUID extracts the UID from the provided Options struct.
// This function attempts to convert the User ID (UID) string from the Options
// struct to an integer. It returns the UID as an integer and any error encountered
// during the conversion process.
//
// Parameters:
// - options: *Options - A pointer to an Options struct containing user information.
//
// Returns:
// - int: The user ID converted to an integer.
// - error: An error object if the conversion fails, nil otherwise.
func parseUID(options *Options) (int, error) {
	uid, err := strconv.Atoi(options.User.Uid)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

// parseGID extracts the GID from the provided Options struct.
// Similar to parseUID, this function converts the Group ID (GID) string from the
// Options struct to an integer. It returns the GID as an integer and any error
// encountered during the conversion.
//
// Parameters:
// - options: *Options - A pointer to an Options struct containing user group information.
//
// Returns:
// - int: The group ID converted to an integer.
// - error: An error object if the conversion fails, nil otherwise.
func parseGID(options *Options) (int, error) {
	gid, err := strconv.Atoi(options.User.Gid)
	if err != nil {
		return 0, err
	}
	return gid, nil
}

// changeOwnership sets the ownership of the specified file or directory.
// This function changes the ownership of the file or directory at the given path
// to the specified user ID (UID) and group ID (GID).
//
// Parameters:
// - path: string - The file or directory path whose ownership is to be changed.
// - uid: int - The user ID to set as owner.
// - gid: int - The group ID to set as owner.
//
// Returns:
// - error: An error object if the ownership change fails, nil otherwise.
func changeOwnership(path string, uid, gid int) error {
	err := os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

// changePermissions sets the permissions of the specified file or directory.
// This function changes the permissions of the file or directory at the given path
// to the permissions specified in the Options struct.
//
// Parameters:
// - path: string - The file or directory path whose permissions are to be changed.
// - perms: *Options - A pointer to an Options struct containing the new permissions.
//
// Returns:
// - error: An error object if the permission change fails, nil otherwise.
func changePermissions(path string, perms *Options) error {
	err := os.Chmod(path, perms.Perms)
	if err != nil {
		return err
	}
	return nil
}
