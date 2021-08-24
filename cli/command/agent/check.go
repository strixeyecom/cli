package agent

import (
	`github.com/spf13/cobra`
	`github.com/spf13/viper`
	
	`github.com/strixeyecom/cli/api/user/agent`
	`github.com/strixeyecom/cli/domain/cli`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	check if current machine supports installing strixeye
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// CheckCommand is a constructor for agent subcommand check.
// It checks for requirements on current machine for installing strixeyed
func CheckCommand() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Command to check install requirements",
		Long: `Command to check install requirements
Every agent has different status. Some are already installed on another machine,
some work with docker. Some wants to use a local database,
some external. This command checks if your current agent is suitable for this machine.
To change your agent, run

strixeye configure agent
`,
		RunE: checkHostCmd,
	}
	
	// declaring local flags used by get trip commands.
	
	return checkCmd
}

// checkHostCmd implements GetCommand logic.
func checkHostCmd(cmd *cobra.Command, _ []string) error {
	var (
		cliConfig cli.Cli
		err       error
	)
	
	// get cli config for authentication
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}
	
	// get agent config from remote.
	agentConfig, err := agent.GetAgentConfig(cliConfig)
	if err != nil {
		return err
		
	}
	
	// check if this host machine supports installing selected agent
	err = agentConfig.CheckIfHostSupports()
	if err != nil {
		return err
	}
	
	return nil
}
