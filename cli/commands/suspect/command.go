package suspect

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	Top level suspicion command.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// NewSuspectCommand serves the base command for querying suspicions.
//
// Adds flags and subcommands for suspicion commands here,
func NewSuspectCommand() *cobra.Command {
	var suspectCommand = &cobra.Command{
		Use:   "suspect",
		Short: "Query and Play with suspects",
		Long:  `Query and play with suspects in the Strixeye Agent of your choice.`,
		RunE:  ShowHelp(os.Stdout),
	}
	
	suspectCommand.AddCommand(
		GetCommand(),
	)
	return suspectCommand
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOut(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
