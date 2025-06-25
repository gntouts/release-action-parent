package main

import (
	"fmt"
	"os"
	"strings"
)

var version string

// processArgs handles the logic for processing command line arguments
func processArgs(args []string) (output string, shouldExit bool, exitCode int) {
	// Handle --version flag
	if len(args) == 1 && args[0] == "--version" {
		versionString := fmt.Sprintf("echo %s", version)
		return versionString, true, 0
	}

	// Handle empty args
	if len(args) == 0 {
		return "", false, 0
	}

	// Echo all arguments
	return strings.Join(args, " "), false, 0
}

// echoOutput handles the actual output printing
func echoOutput(output string) {
	fmt.Println(output)
}

func main() {
	args := os.Args[1:]

	output, shouldExit, exitCode := processArgs(args)

	if shouldExit {
		echoOutput(output)
		os.Exit(exitCode)
		return
	}

	echoOutput(output)
}
