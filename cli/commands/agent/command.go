package agent

import (
	"io"
	"os"
	
	"github.com/spf13/cobra"
	`github.com/strixeyecom/cli/domain/agent`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	top level command for controlling strixeye agent on current machine
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// NewAgentCommand serves the base command for controlling agent on host machine.
//
func NewAgentCommand() *cobra.Command {
	var agentCommands = &cobra.Command{
		Use:   "agent",
		Short: "Control and manage agent on your host machine",
		Long:  `Install, Uninstall, Reset selected agent on current host machine`,
		RunE:  ShowHelp(os.Stdout),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// unfortunately, cobra doesn't support chaining persistent pre run error functions for now.
			// So, we are going to call parent pre run as well
			if parent := cmd.Parent(); parent != nil {
				if parent.PersistentPreRunE != nil {
					err := parent.PersistentPreRunE(parent, args)
					if err != nil {
						return err
					}
				}
			}
			
			// And finally, the pre run function we want
			return agent.IsCorrectUser()
		},
	}
	
	agentCommands.AddCommand(
		CheckCommand(),
		InstallCommand(),
		UninstallCommand(),
	)
	return agentCommands
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOut(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
