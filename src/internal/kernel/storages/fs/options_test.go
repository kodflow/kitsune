package fs

import (
	"os"
	"os/user"
	"testing"

	"github.com/kodflow/kitsune/src/config"
	"github.com/stretchr/testify/assert"
)

func TestResolveFileOptions(t *testing.T) {
	// Test case 1: options is empty
	result := resolveFileOptions()
	expected := defaultFileOptions()
	assert.Equal(t, expected, result, "should return default file options")

	// Test case 2: options[0] is nil
	options := []*Options{nil}
	result = resolveFileOptions(options...)
	expected = defaultFileOptions()
	assert.Equal(t, expected, result, "should return default file options")

	user, _ := user.Current()

	// Test case 3: options[0] is not nil
	options = []*Options{{User: user, Perms: 0755}}
	result = resolveFileOptions(options...)
	expected = options[0]
	assert.Equal(t, expected, result, "should return options[0]")

	// Test case 4: options[0].User is nil
	options = []*Options{{Perms: 0755}}
	result = resolveFileOptions(options...)
	expected = options[0]
	expected.User = config.USER
	assert.Equal(t, expected, result, "should set options[0].User to config.USER")

	// Test case 5: options[0].Perms is 0
	options = []*Options{{User: user}}
	result = resolveFileOptions(options...)
	expected = options[0]
	expected.Perms = 0644
	assert.Equal(t, expected, result, "should set options[0].Perms to 0644")
}
func TestResolveDirectoryOptions(t *testing.T) {
	// Test case 1: options is empty
	result := resolveDirectoryOptions()
	expected := defaultDirectoryOptions()
	assert.Equal(t, expected, result, "should return default directory options")

	// Test case 2: options[0] is nil
	options := []*Options{nil}
	result = resolveDirectoryOptions(options...)
	expected = defaultDirectoryOptions()
	assert.Equal(t, expected, result, "should return default directory options")

	user, _ := user.Current()

	// Test case 3: options[0] is not nil
	options = []*Options{{User: user, Perms: 0755}}
	result = resolveDirectoryOptions(options...)
	expected = options[0]
	assert.Equal(t, expected, result, "should return options[0]")

	// Test case 4: options[0].User is nil
	options = []*Options{{Perms: 0755}}
	result = resolveDirectoryOptions(options...)
	expected = options[0]
	expected.User = config.USER
	assert.Equal(t, expected, result, "should set options[0].User to config.USER")

	// Test case 5: options[0].Perms is 0
	options = []*Options{{User: user}}
	result = resolveDirectoryOptions(options...)
	expected = options[0]
	expected.Perms = 0755
	assert.Equal(t, expected, result, "should set options[0].Perms to 0755")
}
func TestRemovePerms(t *testing.T) {
	// Test case 1: Remove read permission
	options := &Options{Perms: 0644}
	options.RemovePerms(0044)
	expected := &Options{Perms: 0600}
	assert.Equal(t, expected, options, "should remove read permission")

	// Test case 2: Remove write permission
	options = &Options{Perms: 0644}
	options.RemovePerms(0040)
	expected = &Options{Perms: 0604}
	assert.Equal(t, expected, options, "should remove write permission")

	// Test case 3: Remove execute permission
	options = &Options{Perms: 0755}
	options.RemovePerms(0051)
	expected = &Options{Perms: 0704}
	assert.Equal(t, expected, options, "should remove execute permission")

	// Test case 4: Remove multiple permissions
	options = &Options{Perms: 0777}
	options.RemovePerms(0044 | 0022 | 0011)
	expected = &Options{Perms: 0700}
	assert.Equal(t, expected, options, "should remove multiple permissions")

	// Test case 5: Remove no permissions
	options = &Options{Perms: 0644}
	options.RemovePerms(0)
	expected = &Options{Perms: 0644}
	assert.Equal(t, expected, options, "should not remove any permissions")
}

func TestPerms(t *testing.T) {
	defer os.RemoveAll(Kitsune)

	CreateFile(VALID_FILE_PATH)
	u, _ := user.Current()

	// Test case 1: valid options
	options := &Options{
		User: &user.User{
			Uid: u.Uid,
			Gid: u.Gid,
		},
		Perms: 0644,
	}
	err := perms(VALID_FILE_PATH, options)
	assert.NoError(t, err, "should not return an error")

	// Test case 2: invalid user UID
	options = &Options{
		User: &user.User{
			Uid: "invalid",
			Gid: "1000",
		},
		Perms: 0644,
	}
	err = perms(VALID_FILE_PATH, options)
	assert.Error(t, err)

	// Test case 3: invalid user GID
	options = &Options{
		User: &user.User{
			Uid: "1000",
			Gid: "invalid",
		},
		Perms: 0644,
	}
	err = perms(VALID_FILE_PATH, options)
	assert.Error(t, err)

	// Test case 4: failed to change ownership
	options = &Options{
		User: &user.User{
			Uid: "1000",
			Gid: "1000",
		},
		Perms: 0644,
	}
	err = perms(VALID_FILE_PATH, options)
	assert.Error(t, err)

	// Test case 5: failed to change permissions
	options = &Options{
		User: &user.User{
			Uid: "1000",
			Gid: "1000",
		},
		Perms: 0644,
	}
	err = perms(VALID_FILE_PATH, options)
	assert.Error(t, err)

	// Test case 6: failed to change permissions
	options = &Options{
		User: &user.User{
			Uid: "1000",
			Gid: "1000",
		},
		Perms: 0644,
	}

	err = perms(INVALID_FILE_PATH, options)
	assert.Error(t, err)

	// Test case 1: successful permission change
	err = changePermissions(VALID_FILE_PATH, &Options{Perms: 0644})
	assert.NoError(t, err, "should not return an error")

	// Test case 2: failed permission change
	err = changePermissions(INVALID_FILE_PATH, &Options{Perms: 0755})
	assert.Error(t, err)
}
