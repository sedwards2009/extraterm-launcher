package argsparser

import (
	"testing"
)

func TestParseWindow(t *testing.T) {

	testData := []string{
		"extraterm-launcher",
		"-w",
		"1",
	}

	parsedArgs, errorString := Parse(&testData)

	if errorString != nil {
		t.Errorf("Parse error: %s", *errorString)
	}

	if parsedArgs.Window == nil || *parsedArgs.Window != "1" {
		t.Errorf("parsedArgs.Window is bad. Got '%s'", *parsedArgs.Window)
	}
}

func TestParseCommand(t *testing.T) {

	testData := []string{
		"extraterm-launcher",
		"-c",
		"extraterm:window.listAll",
	}

	parsedArgs, errorString := Parse(&testData)

	if errorString != nil {
		t.Errorf("Parse error: %s", *errorString)
	}

	if parsedArgs.CommandName == nil || *parsedArgs.CommandName != "extraterm:window.listAll" {
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

	if errorString != nil {
		t.Errorf("Parse error: %s", *errorString)
	}

	if parsedArgs.CommandName == nil || *parsedArgs.CommandName != "extraterm:window.listAll" {
		t.Errorf("parsedArgs.commandName is bad")
	}

	if _, ok := parsedArgs.CommandParameters["--foo-bar"]; !ok {
		t.Errorf("CommandParameters was missing '--foo-bar'.")
	}
}
