// +build !windows
package settings

import (
	"os"
	"path/filepath"
)

// IpcRunPath returns the absolute path to Extraterm's ipc.run file.
func IpcRunPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".config/extraterm/ipc.run")
}

// ExtratermExeName is the name of the Extraterm main executable.
const ExtratermExeName = "extraterm"
