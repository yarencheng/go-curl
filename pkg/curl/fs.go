package curl

import "github.com/spf13/afero"

// DefaultFs returns the default file system for the current platform.
// In native builds, it returns an OS-backed file system.
func DefaultFs() afero.Fs {
	return afero.NewOsFs()
}
