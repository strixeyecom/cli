package main

import (
	"os"
	
	`github.com/spf13/cobra`
	
	"github.com/strixeyecom/cli/cli/commands/commands"
)

// This is the entrypoint for strixeye cli.
// Constructor for NewStrixeyeCommand is the base command namely "strixeye".
// The constructor adds subcommands for strixeye such as suspect,suspicion,configure.
//
// Usually, configuration initialization takes place in here,
// main package. But I wanted to make that happen in commands package too,
// for the sake of better testability.

var rootCmd *cobra.Command

func main() {
	rootCmd = commands.NewStrixeyeCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
