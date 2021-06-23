package configure

import (
	"github.com/spf13/cobra"
	`github.com/usestrix/cli/cli/commands/repository/ux`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	handles configuring agent for StrixEye CLI
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// NewConfigureAgentCommand constructor for configure agent command.
// If the user is authenticated,
// it can choose agent id's from list of available agents or directly via command line flags.
func NewConfigureAgentCommand() *cobra.Command {
	var configureCommand = &cobra.Command{
		Use:   "agent",
		Short: "Change default agent to work with",
		Long: `A default agent must be selected to work with because StrixEye's CLI is designed to work
with a single StrixEye Agent at a time. `,
		RunE: ux.SetupAgent,
		
	}
	
	// Add subcommands
	configureCommand.AddCommand()
	
	return configureCommand
}
