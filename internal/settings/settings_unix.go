// +build !windows,!darwin

/*
 * Copyright 2021 Simon Edwards <simon@simonzone.com>
 *
 * This source code is licensed under the MIT license which is detailed in the LICENSE.txt file.
 */
package settings

import (
	"os"
	"path/filepath"
)

// IpcRunPath returns the absolute path to Extraterm's ipc.run file.
func IpcRunPath() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgConfigHome) != 0 {
		return filepath.Join(xdgConfigHome, "extraterm/ipc.run")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".config/extraterm/ipc.run")
}

// ExtratermExeName is the name of the Extraterm main executable.
const ExtratermExeName = "extraterm_main"
