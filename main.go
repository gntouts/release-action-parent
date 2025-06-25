package main

import (
	"fmt"
	"os"
	"strings"
)

var version string

func main() {
	args := os.Args[1:]

	// Handle --version flag
	if len(args) == 1 && args[0] == "--version" {
		fmt.Println(version)
		return
	}

	// Echo all arguments
	fmt.Println(strings.Join(args, " "))
}
