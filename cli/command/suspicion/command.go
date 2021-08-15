package suspicion

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

// NewSuspicionCommand serves the base command for querying suspicions.
//
// Adds flags and subcommands for suspicion commands here,
func NewSuspicionCommand() *cobra.Command {
	var suspicionCommand = &cobra.Command{
		Use:   "suspicion",
		Short: "Query and Play with suspicions",
		Long:  `Query and play with suspicions in the Strixeye Agent of your choice.`,
		RunE:  ShowHelp(os.Stdout),
	}

	suspicionCommand.AddCommand(
		GetCommand(),
	)
	return suspicionCommand
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOut(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
