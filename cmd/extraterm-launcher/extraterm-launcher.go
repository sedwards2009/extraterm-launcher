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
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	parsedArgs, errorString := argsparser.Parse(&os.Args)
	if len(errorString) != 0 {
		panic(errorString)
	}

	url := launchMainExecutable()

	exitCode := 0
	if len(parsedArgs.Commands) == 0 {
		runShowWindowCommand(url)
	} else {
		exitCode = runAllCommands(url, parsedArgs)
	}
	os.Stdout.Sync()
	os.Exit(exitCode)
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

func runShowWindowCommand(url string) {
	command := argsparser.MakeCommand()
	showCommandName := string("extraterm:window.show")
	command.CommandName = &showCommandName
	runCommand(url, command)
}

func runAllCommands(url string, parsedArgs *argsparser.CommandLineArguments) int {
	var window *string = nil

	if len(parsedArgs.Commands) != 1 {
		fmt.Println("[")
	}

	for i, command := range parsedArgs.Commands {
		if command.Window != nil {
			window = command.Window
		}

		if command.Window == nil {
			command.Window = window
		}

		httpStatusCode, jsonResult := runCommand(url, command)
		if httpStatusCode < 200 || httpStatusCode >= 300 {
			fmt.Println(jsonResult)
			if len(parsedArgs.Commands) != 1 {
				fmt.Println("[")
			}
			return 1
		} else {

			fmt.Print(jsonResult)
			if i != len(parsedArgs.Commands)-1 {
				fmt.Println(",")
			} else {
				fmt.Println("")
			}
		}
	}

	if len(parsedArgs.Commands) != 1 {
		fmt.Println("[")
	}
	return 0
}

type CommandJson struct {
	CommandName *string           `json:"command"`
	Window      *string           `json:"window"`
	Args        map[string]string `json:"args"`
}

func runCommand(url string, command *argsparser.Command) (httpStatusCode int, jsonBody string) {
	payload := new(CommandJson)
	payload.CommandName = command.CommandName
	payload.Window = command.Window
	payload.Args = wordcase.KababCaseToCamelCaseMapKeys(command.CommandParameters)

	payloadString, _ := json.Marshal(payload)

	buf := strings.NewReader(string(payloadString))
	resp, err := http.Post(url+"/command", "application/json", buf)
	if err != nil {
		panic(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	jsonResult := string(bodyBytes)
	return resp.StatusCode, jsonResult
}
