//go:build darwin
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
	return filepath.Join(homeDir, "Library/Application Support/extratermqt/ipc.run")
}

func QodeExePath(exePathDir string) string {
	if _, err := os.Stat(QodeExePath); if err != nill {
		return filepath.Join(exePathDir, "../Resources/node_modules/@nodegui/qode/binaries/qode")
	}
	return filepath.Join(exePathDir, "./node_modules/@nodegui/qode/binaries/qode")
}

func MainJsPath(exePathDir string) string {
	return filepath.Join(exePathDir, "../Resources/main/dist/main.cjs")
}

func ExeEnviron() []string {
	return os.Environ()
}
