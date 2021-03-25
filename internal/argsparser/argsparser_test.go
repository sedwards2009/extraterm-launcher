/*
 * Copyright 2021 Simon Edwards <simon@simonzone.com>
 *
 * This source code is licensed under the MIT license which is detailed in the LICENSE.txt file.
 */
package argsparser

import (
	"testing"
)

func TestParseWindow(t *testing.T) {

	testData := []string{
		"extraterm-launcher",
		"-c",
		"extraterm:window.show",
		"-w",
		"1",
	}

	parsedArgs, errorString := Parse(&testData)

	if len(errorString) != 0 {
		t.Errorf("Parse error: %s", errorString)
	}

	if parsedArgs.Commands[0].Window == nil || *parsedArgs.Commands[0].Window != "1" {
		t.Errorf("parsedArgs.Window is bad. Got '%s'", *parsedArgs.Commands[0].Window)
	}
}

func TestParseCommand(t *testing.T) {

	testData := []string{
		"extraterm-launcher",
		"-c",
		"extraterm:window.listAll",
	}

	parsedArgs, errorString := Parse(&testData)

	if len(errorString) != 0 {
		t.Errorf("Parse error: %s", errorString)
	}

	if parsedArgs.Commands[0].CommandName == nil || *parsedArgs.Commands[0].CommandName != "extraterm:window.listAll" {
		t.Errorf("parsedArgs.commandName is bad")
	}
}

func TestParseCommandParameters(t *testing.T) {

	testData := []string{
		"extraterm-launcher",
		"-c",
		"extraterm:window.listAll",
		"--foo-bar",
		"true",
	}

	parsedArgs, errorString := Parse(&testData)

	if len(errorString) != 0 {
		t.Errorf("Parse error: %s", errorString)
	}

	if parsedArgs.Commands[0].CommandName == nil || *parsedArgs.Commands[0].CommandName != "extraterm:window.listAll" {
		t.Errorf("parsedArgs.commandName is bad")
	}

	if _, ok := parsedArgs.Commands[0].CommandParameters["--foo-bar"]; !ok {
		t.Errorf("CommandParameters was missing '--foo-bar'.")
	}
}

func TestParseMultipleCommands(t *testing.T) {
	testData := []string{
		"extraterm-launcher",
		"-c",
		"extraterm:window.listAll",
		"-c",
		"extraterm:window.show",
	}

	parsedArgs, errorString := Parse(&testData)

	if len(errorString) != 0 {
		t.Errorf("Parse error: %s", errorString)
	}

	if len(parsedArgs.Commands) != 2 {
		t.Errorf("len(parsedArgs.Commands) is %d, but expected 2.", len(parsedArgs.Commands))
	}
	if *parsedArgs.Commands[0].CommandName != "extraterm:window.listAll" {
		t.Errorf("parsedArgs.Commands[0].CommandName was %s, expected extraterm:window.listAll.",
			*parsedArgs.Commands[0].CommandName)
	}
	if *parsedArgs.Commands[1].CommandName != "extraterm:window.show" {
		t.Errorf("parsedArgs.Commands[1].CommandName was %s, expected extraterm:window.show.",
			*parsedArgs.Commands[1].CommandName)
	}
}

func TestParseBareArgs(t *testing.T) {
	testData := []string{
		"extraterm-launcher",
		"/home/sbe",
	}

	parsedArgs, errorString := Parse(&testData)

	if len(errorString) != 0 {
		t.Errorf("Parse error: %s", errorString)
	}
	if len(parsedArgs.BareArgs) != 1 {
		t.Errorf("parsedArgs.BareArgs != 1")
	}
}
