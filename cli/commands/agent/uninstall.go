package agent

import (
	`fmt`
	`os`
	`path/filepath`
	`strings`
	
	"github.com/fatih/color"
	`github.com/manifoldco/promptui`
	`github.com/pkg/errors`
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	agent2 "github.com/usestrix/cli/domain/agent"
	"github.com/usestrix/cli/domain/consts"
	
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
	
	// Uninstalling a new agent while one is still running is a bad practice in our current system
	err = agent2.CheckIfAnotherAgentRunning()
	if err != nil {
		return err
	}
	
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
	
	err = agent2.StopDaemon()
	if err != nil {
		return err
	}
	color.Red("Stopped StrixEye Daemon")
	
	err = os.Remove(filepath.Join(consts.DaemonDir, consts.DaemonName))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	color.Red("Removed StrixEye Daemon")
	
	err = os.Remove(filepath.Join(consts.ServiceDir, consts.ServiceFile))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	color.Red("Removed StrixEye Service Files")
	
	// Remove StrixEye Volumes
	prompt := promptui.Prompt{
		Label:     "Would you like to remove all StrixEye Agent related data? (recommended)",
		IsConfirm: true,
	}
	
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Not removing volumes\n")
		return nil
	}
	
	if strings.EqualFold(result, "y") {
		err = prune(agentConfig.Config.Deployment)
		if err != nil {
			return err
		}
	}
	
	color.Red("Uninstall completed successfully")
	return nil
}

// prune removes all stored volume data on host machine depending on deployment type
func prune(deploymentName string) error {
	var (
		err error
	)
	
	// 	remove kubernetes volumes and networks
	if deploymentName == consts.KubernetesDeployment {
	
	}
	
	// 	remove docker volumes and networks
	if deploymentName == consts.DockerDeployment {
		err = agent2.RemoveDockerVolumeByName(agent2.DockerBrokerVolumeName)
		if err != nil {
			return err
		}
		
		err = agent2.RemoveDockerVolumeByName(agent2.DockerDatabaseVolumeName)
		if err != nil {
			return err
		}
		return nil
	}
	
	return errors.Errorf("no such strixeye volume :%s", deploymentName)
}
