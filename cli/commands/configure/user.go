package configure

import (
	"github.com/spf13/cobra"
	`github.com/usestrix/cli/cli/commands/repository/ux`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	handles configuring user for StrixEye CLI
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// NewConfigureUserCommand returns the command for
//  strixeye configure user
// Which lets users authenticate via prompts or flags
func NewConfigureUserCommand() *cobra.Command {
	// serveCmd represents the request command
	var configureUserCommand = &cobra.Command{
		Use:   "user",
		Short: "Change default user to work with",
		Long: `A default user must be selected to work with because StrixEye's CLI is designed to work
with a single StrixEye user's API Token.

Authentication will be executed via this user's API Token' `,
		RunE: ux.SetupUser,
	}
	
	// Add subcommands
	configureUserCommand.AddCommand()
	
	return configureUserCommand
}



