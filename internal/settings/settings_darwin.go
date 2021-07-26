// +build darwin

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, "Library/Application Support/extraterm-qt/ipc.run")
}

const QodeExePath = "node_modules/@nodegui/qode/binaries/qode"
const MainJsPath = "main/dist/main.js"
