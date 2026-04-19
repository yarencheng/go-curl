//go:build js && wasm

package curl

import "github.com/spf13/afero"

// DefaultFs returns the default file system for the current platform.
// In WASM builds, it returns an in-memory file system.
func DefaultFs() afero.Fs {
	return afero.NewMemMapFs()
}
