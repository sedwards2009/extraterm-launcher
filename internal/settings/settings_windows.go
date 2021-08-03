// +build windows

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

func IpcRunPath() string {
	appData := os.Getenv("APPDATA")
	return filepath.Join(appData, "extraterm-qt", "ipc.run")
}

const QodeExePath = "qode"
const MainJsPath = "main/dist/main.js"
