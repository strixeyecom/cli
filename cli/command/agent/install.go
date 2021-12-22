package agent

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/strixeyecom/cli/cli/command/repository/ux"
	agent2 "github.com/strixeyecom/cli/domain/agent"
	"github.com/strixeyecom/cli/domain/consts"

	"github.com/strixeyecom/cli/api/user/agent"
	"github.com/strixeyecom/cli/domain/cli"
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

	checkCmd.Flags().Bool(
		"interactive", false, "--interactive if you want to configure StrixEye CLI during installation",
	)

	return checkCmd
}

func getCredentials(cmd *cobra.Command) (cli.Cli, error) {
	var (
		cliConfig cli.Cli
		err       error
	)

	isInteractive, err := cmd.Flags().GetBool("interactive")
	if err != nil {
		return cli.Cli{}, err
	}

	if !isInteractive {
		// get cli config for authentication
		err = viper.Unmarshal(&cliConfig)
		if err != nil {
			return cliConfig, err
		}

		return cliConfig, nil
	}

	// if interactive install,
	err = ux.SetupUser(cmd, nil)
	if err != nil {
		return cli.Cli{}, err
	}

	// if successfully setup user, than we can set up the agent
	err = ux.SetupAgent(cmd, nil)
	if err != nil {
		return cli.Cli{}, err
	}
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

	// Installing a new agent while one is still running is a bad practice in our current system
	err = agent2.CheckIfAnotherAgentRunning()
	if err != nil {
		return err
	}

	cliConfig, err = getCredentials(cmd)
	if err != nil {
		return err
	}

	// early cut if bad credentials
	if cliConfig.UserAPIToken == "" {
		return errors.Errorf(
			`empty user api token during installation. Please check out documentation.
 If you haven't set up your cli, you can set it up during installation:

	$ strixeye agent install --interactive
`,
		)
	}
	// early cut if bad agent id
	if cliConfig.AgentID == "" {
		return errors.Errorf(
			`empty agent id in configuration. Please check out documentation.
 If you haven't set up your cli, you can set it up during installation:

	$ strixeye agent install --interactive
`,
		)
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

	err = loginToDocker(agentConfig, cliConfig)
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
	err = DownloadDaemonBinary(
		cliConfig.UserAPIToken, agentConfig.Token, versions.Manager,
		cliConfig.DownloadDomain,
	)

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
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(consts.WorkingDir, 0600)
		if err != nil {
			return err
		}
	}

	// 	create config directory
	_, err = os.Stat(consts.ConfigDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		err = os.Mkdir(consts.ConfigDir, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// loginToDocker is necessary because StrixEye images are kept in a secure registry.
func loginToDocker(agentConfig agent2.AgentInformation, cliConfig cli.Cli) error {
	testRefStr := fmt.Sprintf("%s/%s", cliConfig.DockerRegistry, "hello-world")

	cmd := exec.CommandContext(context.Background(), "docker", "pull", testRefStr)

	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	err := cmd.Run()

	if err != nil {
		if agentConfig.Config.Deployment == "docker" {
			color.Red(
				`
Docker can not fetch test image from StrixEye registry.
If you haven't logged in to Strixeye Docker Registry yet, you need to login to StrixEye registry at %s.

Here is a direct command you can execute:

$ sudo docker login --username strixeye --password %s %s

If you don't want to show credentials in your history, following is a custom command you can use.
Logging in to a docker registry is a basic procedure that is documented by docker as well. https://docs.docker.com/engine/reference/commandline/login/

$ strixeye inspect user_api_token | sudo docker login --username strixeye --password-stdin $(
strixeye inspect docker_registry)
`, cliConfig.DockerRegistry, agentConfig.Token, cliConfig.DockerRegistry,
			)
		} else if agentConfig.Config.Deployment == "kubernetes" {
			color.Red(
				`
Kubernetes can not fetch test image from StrixEye registry.
If you haven't logged in to Strixeye Container Registry yet, you need to login to StrixEye registry at %s.

Here is a direct command you can execute:

$ kubectl create secret docker-registry strixeye-cred --docker-server=$(strixeye inspect docker_registry) --docker-username=strixeye --docker-password=$(strixeye inspect user_api_token)
`,
			)
		}

		return err
	}

	return nil
}
