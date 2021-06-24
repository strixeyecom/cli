package agent

import (
	`os`
	`path/filepath`
	
	`github.com/fatih/color`
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	agent2 `github.com/usestrix/cli/domain/agent`
	`github.com/usestrix/cli/domain/consts`
	
	"github.com/usestrix/cli/api/user/agent"
	"github.com/usestrix/cli/domain/cli"
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

// UninstallCommand is a constructor for agent subcommand check.
// It checks for requirements on current machine for installing strixeyed
func UninstallCommand() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Command to uninstall agent",
		Long: `Command to uninstall agent

This command will uninstall your current agent on this machine.
You will need root permissions for this command.
To change your agent, run
strixeye configure agent
`,
		RunE: uninstallAgentCmd,
	}
	
	// declaring local flags used by get trip commands.
	
	return checkCmd
}

// installAgentCmd implements GetCommand logic.
func uninstallAgentCmd(cmd *cobra.Command, _ []string) error {
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
	_, err = agent.GetAgentConfig(cliConfig)
	if err != nil {
		return err
		
	}
	
	err = agent2.StopDaemon()
	if err != nil {
		return err
	}
	color.Red("Stopped StrixEye Daemon")
	
	err = os.Remove(filepath.Join(consts.DaemonDir, consts.DaemonName))
	if err != nil {
		return err
	}
	color.Red("Removed StrixEye Daemon")
	
	
	return nil
}
