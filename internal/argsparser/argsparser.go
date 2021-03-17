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

const commandFlagStr = "-c"
const windowFlagStr = "-w"

func Parse(args *[]string) (parsed *CommandLineArguments, errString string) {
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
			case windowFlagStr:
				return result, "-w parameter must be used after -c."

			case commandFlagStr:
				state = CommandName

			default:
				return result, fmt.Sprintf("Unknown parameter '%s'.", item)
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
			if item == commandFlagStr {
				if command.CommandName != nil {
					result.Commands = append(result.Commands, command)
					command = MakeCommand()
				}
				state = CommandName

			} else if item == windowFlagStr {
				state = WindowValue

			} else if strings.HasPrefix(item, "--") {
				commandFlag = item
				state = CommandValue
			} else {
				return result, fmt.Sprintf("Parameters to commands must start with '--'. Found: '%s'.", item)
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
		return result, fmt.Sprintf("Command parameter '%s' does have a valkue.", commandFlag)
	}

	if state == WindowValue {
		return result, "No value for window option was given."
	}

	return result, ""
}
