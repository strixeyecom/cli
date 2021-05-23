package main

import (
	"os"

	"github.com/usestrix/cli/cli/commands/commands"
)

// This is the entrypoint for strixeye cli.
// Constructor for NewStrixeyeCommand is the base command namely "strixeye".
// The constructor adds subcommands for strixeye such as suspect,suspicion,configure.
//
// Usually, configuration initialization takes place in here,
// main package. But I wanted to make that happen in commands package too,
// for the sake of better testability.
func main() {
	cmd := commands.NewStrixeyeCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
