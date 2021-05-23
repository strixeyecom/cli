package configure

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/usestrix/cli/api/user"
	"github.com/usestrix/cli/domain/config"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// setting local variable definitions
			var (
				err           error
				cliConfig     config.Cli
				selectedAgent config.AgentInformation
				selectedID    string
			)

			// try to get from flags.
			selectedID, err = cmd.Flags().GetString("set-agent-id")
			if err != nil {
				return err
			}

			// if flags aren't provided, let user choose from a select list.
			if selectedID == "" {
				// get cli config if exists.
				err = viper.Unmarshal(&cliConfig)
				if err != nil {
					return err
				}

				// select agent with current cli configuration
				selectedAgent, err = selectAgent(cliConfig)
				if err != nil {
					return err
				}

				selectedID = selectedAgent.ID
			}

			// print selected information
			color.Blue("Selected agent: %s", selectedID)

			// store selected agent to configuration file.
			viper.Set("CURRENT_AGENT_ID", selectedID)
			err = viper.WriteConfig()
			if err != nil {
				return errors.Wrap(err, "failed to save config. Do you have permissions on the filesystem? ")
			}

			return nil
		},
	}

	// Add subcommands
	configureCommand.AddCommand()

	// Add local flags
	configureCommand.Flags().StringP(
		"set-agent-id", "s", "",
		`Set agent directly without using select functionality`,
	)

	return configureCommand
}

// selectAgent displays a select list of agents retrieved from user api.
func selectAgent(cliConfig config.Cli) (config.AgentInformation, error) {
	var (
		err           error
		agents        []config.AgentInformation
		selectedAgent config.AgentInformation
	)

	// fetch agents from user api
	agents, err = user.GetAgents(cliConfig)
	if err != nil {
		return config.AgentInformation{}, err
	}

	// prepare inputs. create a string slice of agents.
	var selectedIndex int
	promptItems := make([]string, len(agents))

	for i := range agents {
		promptItems[i] = agents[i].String()
	}

	question := []*survey.Question{
		{
			Name: "agent",
			Prompt: &survey.Select{
				Message: "Choose an agent:",
				Options: promptItems,
			},
		},
	}

	// perform the questions
	err = survey.Ask(question, &selectedIndex)
	if err != nil {
		return config.AgentInformation{}, err
	}

	selectedAgent = agents[selectedIndex]

	return selectedAgent, nil
}
