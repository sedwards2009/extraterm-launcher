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
	if len(parsedArgs.BareArgs) > 0 {
		exitCode = runOpenWindowAtCommand(url, parsedArgs.BareArgs[0])
		if exitCode == 0 {
			exitCode = runShowWindowCommand(url)
		}
	} else if len(parsedArgs.Commands) == 0 {
		cwd, _ := os.Getwd()
		exitCode = runOpenWindowAtCommand(url, cwd)
		if exitCode == 0 {
			exitCode = runShowWindowCommand(url)
		}
	} else {
		exitCode = runAllCommands(url, parsedArgs)
	}
	os.Stdout.Sync()
	os.Exit(exitCode)
}

// Returns: URL to Extraterm's localhost IPC server
func launchMainExecutable() string {
	url := readIpcRunFile(settings.IpcRunPath())
	if url != "" {
		if ping(url) {
			return url
		}
	}

	return runMainExecutable()
}

func readIpcRunFile(ipcRunPath string) string {
	contentBytes, err := os.ReadFile(settings.IpcRunPath())
	if err != nil {
		return ""
	}
	content := string(contentBytes)
	parts := strings.Split(content, "\n")
	_, err = strconv.Atoi(parts[0])
	if err != nil {
		return ""
	}
	return parts[1]
}

// Returns: URL to Extraterm's localhost IPC server
func runMainExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePathDir := filepath.Dir(exePath)
	qodeExePath := settings.QodeExePath(exePathDir)
	mainJsPath := settings.MainJsPath(exePathDir)

	cmd := exec.Command(qodeExePath, mainJsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = settings.ExeEnviron()
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start the main Extraterm executable. %s\n", err)
		panic(nil)
	}

	return waitForMainExecutableToAppear(settings.IpcRunPath())
}

// Returns: URL to Extraterm's localhost IPC server or empty string if Extraterm didn't respond.
func waitForMainExecutableToAppear(ipcRunPath string) string {
	sleepTime := 250 * time.Millisecond
	retryTime := 10 * time.Second

	for i := time.Duration(0); i < retryTime; i += sleepTime {
		url := readIpcRunFile(ipcRunPath)
		if ping(url) {
			return url
		}

		time.Sleep(sleepTime)
	}

	return "" // Time-out
}

func runShowWindowCommand(url string) int {
	command := argsparser.MakeCommand()
	showCommandName := string("extraterm:window.showAll")
	command.CommandName = &showCommandName

	httpStatusCode, jsonResult := runCommand(url, command)
	if isErrorHttpStatusCode(httpStatusCode) {
		fmt.Println(jsonResult)
		return 1
	}
	return 0
}

func runOpenWindowAtCommand(url string, path string) int {
	command := argsparser.MakeCommand()
	newTerminalCommandName := string("extraterm:window.newTerminal")
	command.CommandName = &newTerminalCommandName
	command.CommandParameters["--working-directory"] = path

	httpStatusCode, jsonResult := runCommand(url, command)
	if isErrorHttpStatusCode(httpStatusCode) {
		fmt.Println(jsonResult)
		return 1
	}
	return 0
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
		if isErrorHttpStatusCode(httpStatusCode) {
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

func isErrorHttpStatusCode(httpStatusCode int) bool {
	return httpStatusCode < 200 || httpStatusCode >= 300
}

func ping(url string) bool {
	resp, err := http.Get(url + "/ping")
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	message := string(bodyBytes)
	return message == "pong"
}
