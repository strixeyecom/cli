package agent

import (
	"fmt"
	"os"
	
	"github.com/pkg/errors"
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

// InstallCommand is a constructor for agent subcommand check.
// It checks for requirements on current machine for installing strixeyed
func InstallCommand() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "install",
		Short: "Command to install agent",
		Long: `Command to install agent

This command will install your current agent on this machine.
You will need root permissions for this command.
To change your agent, run
strixeye configure agent
`,
		RunE: installAgentCmd,
	}
	
	// declaring local flags used by get trip commands.
	checkCmd.Flags().String(
		"user-api-token", "", "--token {your StrixEye User API Token}",
	)
	
	checkCmd.Flags().String(
		"agent-id", "", "--agent-id {StrixEye Agent ID of your choice}",
	)
	
	return checkCmd
}

func getCredentials(cmd *cobra.Command) (cli.Cli, error) {
	var (
		cliConfig cli.Cli
		err       error
	)
	
	// get cli config for authentication
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return cliConfig, err
	}
	
	return cliConfig, nil
}

// installAgentCmd implements GetCommand logic.
func installAgentCmd(cmd *cobra.Command, _ []string) error {
	var (
		cliConfig cli.Cli
		err       error
	)
	
	
	cliConfig, err = getCredentials(cmd)
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
	
	// get latest versions
	versions, err := agent.GetLatestVersions(cliConfig)
	if err != nil {
		return err
	}
	
	// create necessary directories and files.
	err = createPaths(agentConfig)
	if errors.Is(err, os.ErrExist) {
		fmt.Println(err)
	} else if err != nil {
		return err
	}
	
	// create service file depending on os/arch and deployment type
	err = agentConfig.CreateServiceFile()
	if err != nil {
		return err
	}
	
	// Save agent config file
	a := agent2.Agent{
		Versions: versions, Auth: agent2.Auth{
			AgentID:    agentConfig.ID,
			AgentToken: agentConfig.Token,
		},
		Addresses: agentConfig.Config.Addresses,
	}
	err = agent2.SaveAgentConfig(a)
	if err != nil {
		return err
	}
	
	fmt.Println("Starting download process.")
	// download tarball, decompress and place the binary
	err = DownloadDaemonBinary(cliConfig.UserAPIToken, agentConfig.Token, versions.Manager)
	if err != nil {
		return err
	}
	agent2.InstallCompleted()
	return nil
}

// createPaths creates paths depending on the host machine os/arch. For example,
// working directory in win and unix are different, decided on the compile time.
func createPaths(agentInformation agent2.AgentInformation) error {
	// create working directory
	_, err := os.Stat(consts.WorkingDir)
	switch {
	case errors.Is(os.ErrNotExist, err):
		err = os.Mkdir(consts.WorkingDir, 0600)
		if err != nil {
			return err
		}
	case err == nil:
		return errors.WithMessagef(os.ErrExist, "directory %s exists", consts.WorkingDir)
	default:
		return errors.WithMessagef(err, "can not create directory %s", consts.WorkingDir)
	}
	
	// 	create config directory
	_, err = os.Stat(consts.ConfigDir)
	switch {
	case errors.Is(os.ErrNotExist, err):
		err = os.Mkdir(consts.ConfigDir, 0600)
		if err != nil {
			return err
		}
	case err == nil:
		return errors.WithMessagef(os.ErrExist, "directory %s exists", consts.ConfigDir)
	default:
		return errors.WithMessagef(err, "can not create directory %s", consts.ConfigDir)
	}
	
	return nil
}
