package configure

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

/*
	Created by aomerk at 5/22/21 for project cli
*/

/*
	This is the base group for configure commands and doesn't do much except printing.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// NewConfigureCommand constructor for configure command.
// This is the base group for configure commands and doesn't do much except printing.
func NewConfigureCommand() *cobra.Command {
	// serveCmd represents the request command
	var configureCommand = &cobra.Command{
		Use:   "configure",
		Short: "Setup your cli configuration and credentials.",
		Long: `Setup your cli configuration and credentials.
At bare, you will need a User API Token from StrixEye and a chosen StrixEye Agent's id.`,
		RunE: ShowHelp(os.Stdout),
	}

	// Add subcommands
	configureCommand.AddCommand(
		NewConfigureAgentCommand(),
		NewConfigureUserCommand(),
		NewInspectCommand(),
	)
	return configureCommand
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOut(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
