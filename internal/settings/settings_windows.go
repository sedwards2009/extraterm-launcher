// +build windows

/*
 * Copyright 2021 Simon Edwards <simon@simonzone.com>
 *
 * This source code is licensed under the MIT license which is detailed in the LICENSE.txt file.
 */
package settings

import (
	"os"
)

func IpcRunPath() string {
	appData := os.Getenv("APPDATA")
	return filepath.Join(appData, ".config/extraterm/ipc.run")
}

const ExtratermExeName = "extraterm.exe"
