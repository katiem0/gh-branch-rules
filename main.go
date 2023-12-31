package main

import (
	"os"

	"github.com/katiem0/gh-branch-rules/cmd"
)

func main() {
	// Instantiate and execute root command
	cmd := cmd.NewCmdRoot()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
