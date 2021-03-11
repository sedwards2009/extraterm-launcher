package argsparser

import (
	"fmt"
	"strings"
)

type CommandLineArguments struct {
	Window            *string
	CommandName       *string
	CommandParameters map[string]string
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
	result := &CommandLineArguments{CommandParameters: map[string]string{}}

	state := ProgramName
	commandFlag := ""

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
			result.Window = &window
			state = Flag

		case CommandName:
			commandName := item
			result.CommandName = &commandName
			state = CommandFlag

		case CommandFlag:
			if strings.HasPrefix(item, "--") {
				commandFlag = item
				state = CommandValue
			} else {
				errorMessage := fmt.Sprintf("Parameters to commands must start with '--'. Found: '%s'.", item)
				return result, &errorMessage
			}
		case CommandValue:
			commandValue := item
			result.CommandParameters[commandFlag] = commandValue
			state = CommandFlag
		}
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
