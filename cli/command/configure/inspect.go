package configure

import (
	"encoding/json"

	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/strixeyecom/cli/api/user/agent"
	agent2 "github.com/strixeyecom/cli/domain/agent"
	"github.com/strixeyecom/cli/domain/cli"
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

// NewInspectCommand returns the command for
//  strixeye configure inspect
// configurations.
//
// Naked run returns cli config, if you pass agent id as parameter, it returns agent configuration.
func NewInspectCommand() *cobra.Command {
	// serveCmd represents the request command
	var configureUserCommand = &cobra.Command{
		Use:   "inspect",
		Short: "Inspect configurations",
		Long: `A default user must be selected to work with because StrixEye's CLI is designed to work
with a single StrixEye user's API Token.

Authentication will be executed via this user's API Token

Naked run returns current StrixEye CLI configuration that you are running right now.

Passing Agent ID as parameters show agent configurations
' `,
		RunE: inspectCmd,
	}
	
	return configureUserCommand
}

// inspectCmd implements command functionality.
func inspectCmd(cmd *cobra.Command, args []string) error {
	var (
		err               error
		AgentInformations = make([]agent2.AgentInformation, len(args))
		cliConfig         cli.Cli
	)

	// get cli config
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}

	// if no args, show cli config
	if len(args) == 0 {
		_, err = pp.Print(cliConfig)
		if err != nil {
			return err
		}
	} else {
		// Show all requested agent configurations from user api
		tmp := cliConfig
		for i, arg := range args {
			tmp.AgentID = arg
			AgentInformations[i], err = agent.GetAgentConfig(cliConfig)
			if err != nil {
				return err
			}
		}
		data, err := json.MarshalIndent(AgentInformations, "", "\t")
		if err != nil {
			return err
		}

		cmd.Println(string(data))
	}

	return nil
}
