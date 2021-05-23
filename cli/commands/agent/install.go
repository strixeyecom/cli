package agent

import (
	`fmt`
	`io`
	`net/http`
	`os`
	`time`
	
	`github.com/fatih/color`
	`github.com/spf13/cobra`
	`github.com/spf13/viper`
	
	`github.com/usestrix/cli/api/user/agent`
	`github.com/usestrix/cli/domain/config`
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
	
	return checkCmd
}

// installAgentCmd implements GetCommand logic.
func installAgentCmd(cmd *cobra.Command, _ []string) error {
	var (
		cliConfig config.Cli
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
	
	// create http request
	installerURL := fmt.Sprintf("https://dashboard.***REMOVED***/download/%s", cliConfig.UserAPIToken)
	client := http.Client{Timeout: time.Second * 10}
	request, err := http.NewRequest(http.MethodGet, installerURL, nil)
	if err != nil {
		return err
	}
	request.Header.Add("accept", "application/json")
	
	// notify user
	color.Blue("Downloading install script from %s", installerURL)
	
	// download installer script
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// open a temporary file to store installer script
	tmpFileName := fmt.Sprintf("/tmp/installer_%d.sh", time.Now().UnixNano())
	f, err := os.Create(tmpFileName)
	if err != nil {
		return err
	}
	
	// write to temporary file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	
	color.Red(
		`!!ATTENTION!!
PLEASE RUN THE FOLLOWING SCRIPT
%s
`, tmpFileName,
	)
	
	return nil
}
