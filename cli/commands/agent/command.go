package agent

import (
	`io`
	`os`
	
	`github.com/spf13/cobra`
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
		Short: "Control and manage agent",
		Long:  `Install, Uninstall, Reset selected agent on current host machine`,
		RunE:  ShowHelp(os.Stdout),
	}
	
	agentCommands.AddCommand(
		CheckCommand(),
		InstallCommand(),
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
