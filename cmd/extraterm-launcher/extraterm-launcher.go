/*
 * Copyright 2021 Simon Edwards <simon@simonzone.com>
 *
 * This source code is licensed under the MIT license which is detailed in the LICENSE.txt file.
 */
package main

import (
	"encoding/json"
	"extraterm-launcher/internal/argsparser"
	"extraterm-launcher/internal/settings"
	"extraterm-launcher/internal/wordcase"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Extraterm launcher")

	parsedArgs, errorString := argsparser.Parse(&os.Args)
	if errorString != nil {
		panic(errorString)
	}

	url := launchMainExecutable()
	fmt.Printf("Main executable launched and base URL is %s\n", url)

	if len(parsedArgs.Commands) == 0 {
		command := argsparser.MakeCommand()
		showCommandName := string("extraterm:window.show")
		command.CommandName = &showCommandName
		runCommand(url, command)
	} else {

		var window *string = nil
		for _, command := range parsedArgs.Commands {
			if command.Window != nil {
				window = command.Window
			}

			if command.Window == nil {
				command.Window = window
			}

			runCommand(url, command)
		}
	}
}

func launchMainExecutable() string {
	pid, url := readIpcRunFile(settings.IpcRunPath())

	if pid < 0 {
		url = runMainExecutable()
	}

	return url
}

func readIpcRunFile(ipcRunPath string) (pid int, url string) {
	contentBytes, err := os.ReadFile(settings.IpcRunPath())
	if err != nil {
		pid = -1
		return
	}
	content := string(contentBytes)
	parts := strings.Split(content, "\n")
	pid, err = strconv.Atoi(parts[0])
	if err != nil {
		pid = -1
		return
	}
	url = parts[1]
	return
}

func runMainExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	mainExePath := filepath.Join(filepath.Dir(exePath), settings.ExtratermExeName)
	fmt.Printf("Main executable path %s\n", mainExePath)
	cmd := exec.Command(mainExePath)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start the main Extraterm executable. %s\n", err)
		panic(nil)
	}

	return waitForMainExecutableToAppear(cmd.Process.Pid, settings.IpcRunPath())
}

func waitForMainExecutableToAppear(pid int, ipcRunPath string) string {
	for true {
		filePid, url := readIpcRunFile(ipcRunPath)
		if filePid == pid {
			return url
		}

		_, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to start the main Extraterm executable. It appears to have died.")
			panic(nil)
		}

		time.Sleep(250 * time.Millisecond)
	}

	return "" // Unreachable
}

type CommandJson struct {
	CommandName *string           `json:"command"`
	Window      *string           `json:"window"`
	Args        map[string]string `json:"args"`
}

func runCommand(url string, command *argsparser.Command) {
	payload := new(CommandJson)
	payload.CommandName = command.CommandName
	payload.Window = command.Window
	payload.Args = wordcase.KababCaseToCamelCaseMapKeys(command.CommandParameters)

	payloadString, _ := json.Marshal(payload)
	fmt.Printf("  %s\n", string(payloadString))

	buf := strings.NewReader(string(payloadString))
	resp, err := http.Post(url+"/command", "application/json", buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %d", resp.StatusCode)
}
