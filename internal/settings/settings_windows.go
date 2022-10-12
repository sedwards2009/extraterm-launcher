//go:build windows
// +build windows

/*
 * Copyright 2022 Simon Edwards <simon@simonzone.com>
 *
 * This source code is licensed under the MIT license which is detailed in the LICENSE.txt file.
 */
package settings

import (
	"os"
	"path/filepath"
	"strings"
)

func IpcRunPath() string {
	appData := os.Getenv("APPDATA")
	return filepath.Join(appData, "extratermqt", "Config", "ipc.run")
}

const QodeExePath = "./node_modules/@nodegui/qode/binaries/qode.exe"
const MainJsPath = "main/dist/main.cjs"

func ExeEnviron() []string {
	env := os.Environ()
	exePath, err := os.Executable()
	if err != nil {
		return env
	}

	exePath = filepath.Dir(exePath)

	var extraPath strings.Builder
	dllPaths := []string{
		"main\\resources\\list-fonts-json-binary\\win32-x64",
		"node_modules\\@nodegui\\nodegui\\miniqt\\5.15.2\\msvc2019_64\\bin",
		"node_modules\\@nodegui\\nodegui\\miniqt\\5.15.2\\msvc2019_64\\plugins\\iconengines",
		"node_modules\\@nodegui\\nodegui\\miniqt\\5.15.2\\msvc2019_64\\plugins\\imageformats",
		"node_modules\\@nodegui\\nodegui\\miniqt\\5.15.2\\msvc2019_64\\plugins\\platforms",
		"node_modules\\@nodegui\\nodegui\\miniqt\\5.15.2\\msvc2019_64\\plugins\\styles",
	}

	extraPath.WriteString("Path=")
	for _, dllPath := range dllPaths {
		extraPath.WriteString(filepath.Join(exePath, dllPath))
		extraPath.WriteString(";")
	}

	for _, pair := range env {
		if strings.HasPrefix(pair, "Path=") {
			extraPath.WriteString(pair[5:])
		}
	}

	return append(env, extraPath.String())
}
