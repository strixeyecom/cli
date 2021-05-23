package trip

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	Top level trip command.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// NewTripCommand serves the base command for querying trips.
//
// Adds flags and subcommands for trip commands here,
func NewTripCommand() *cobra.Command {
	var tripCommands = &cobra.Command{
		Use:   "trip",
		Short: "Query and Play with trips",
		Long:  `Query and play with trips in the Strixeye Agent of your choice.`,
		RunE:  ShowHelp(os.Stdout),
	}

	tripCommands.AddCommand(
		GetCommand(),
	)
	return tripCommands
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOut(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
