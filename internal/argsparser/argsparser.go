/*
 * Copyright 2021 Simon Edwards <simon@simonzone.com>
 *
 * This source code is licensed under the MIT license which is detailed in the LICENSE.txt file.
 */
package argsparser

import (
	"fmt"
	"strings"
)

type Command struct {
	Window            *string
	CommandName       *string
	CommandParameters map[string]string
}

type CommandLineArguments struct {
	Commands []*Command
}

func MakeCommand() *Command {
	return &Command{CommandParameters: map[string]string{}}
}

type ParseState int

const (
	ProgramName ParseState = iota
	Flag
	WindowValue
	CommandName
	CommandFlag
	CommandValue
)

func Parse(args *[]string) (parsed *CommandLineArguments, errString *string) {
	result := &CommandLineArguments{Commands: []*Command{}}

	state := ProgramName
	commandFlag := ""

	command := MakeCommand()

	for _, item := range *args {
		switch state {

		case ProgramName:
			state = Flag

		case Flag:
			switch item {

			case "-w", "--window":
				state = WindowValue

			case "-c", "--command":
				state = CommandName

			default:
				errorMessage := fmt.Sprintf("Unknown command line parameter '%s'.", item)
				return result, &errorMessage
			}

		case WindowValue:
			window := item
			command.Window = &window
			state = Flag

		case CommandName:
			commandName := item
			command.CommandName = &commandName
			state = CommandFlag

		case CommandFlag:
			if item == "--" {
				if command.CommandName != nil {
					result.Commands = append(result.Commands, command)
					command = MakeCommand()
				}
				state = Flag
			} else if strings.HasPrefix(item, "--") {
				commandFlag = item
				state = CommandValue
			} else {
				errorMessage := fmt.Sprintf("Parameters to commands must start with '--'. Found: '%s'.", item)
				return result, &errorMessage
			}

		case CommandValue:
			commandValue := item
			command.CommandParameters[commandFlag] = commandValue
			state = CommandFlag
		}
	}

	if command.CommandName != nil {
		result.Commands = append(result.Commands, command)
	}

	if state == CommandValue {
		errorMessage := fmt.Sprintf("Command parameter '%s' does have a valkue.", commandFlag)
		return result, &errorMessage
	}

	if state == WindowValue {
		errorMessage := "No value for window option was given."
		return result, &errorMessage
	}

	return result, nil
}
