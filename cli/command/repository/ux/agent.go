package ux

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/strixeyecom/cli/api/user/agent"
	agent2 "github.com/strixeyecom/cli/domain/agent"
	"github.com/strixeyecom/cli/domain/cli"
)

/*
	Created by aomerk at 6/23/21 for project cli
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// SelectAgentFromList displays a select list of agents retrieved from user api.
func SelectAgentFromList(agents []agent2.AgentInformation) (agent2.AgentInformation, error) {
	var (
		err           error
		selectedAgent agent2.AgentInformation
	)

	// prepare inputs. create a string slice of agents.
	var selectedIndex int
	promptItems := make([]string, len(agents))

	for i := range agents {
		promptItems[i] = agents[i].ID
	}

	question := promptui.Select{
		Label: "Choose an agent:",
		Items: promptItems,
		Size:  45,
	}

	// perform the questions
	selectedIndex, _, err = question.Run()
	if err != nil {
		return agent2.AgentInformation{}, err
	}

	selectedAgent = agents[selectedIndex]

	return selectedAgent, nil
}

func SetupAgent(cmd *cobra.Command, _ []string) error {
	// setting local variable definitions
	var (
		err           error
		cliConfig     cli.Cli
		selectedAgent agent2.AgentInformation
		selectedID    string
	)
	// get cli config if exists.
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}

	// get the list of authorized agents by this user

	// get a list of agents
	agents, err := agent.GetAgents(cliConfig)
	if err != nil {
		return err
	}

	// check if user want to use it directly
	selectedID, err = cmd.Flags().GetString("agent-id")
	if err != nil {
		return err
	}

	// if flags aren't provided, let user choose from a select list.
	isInteractive, err := cmd.Flags().GetBool("interactive")
	if err != nil {
		return err
	}

	// handle interactive setup
	if selectedID == "" || isInteractive {
		selectedID, err = handleInteractive(agents, cliConfig)
		if err != nil {
			return err
		}
	}

	// find entered agent in this user's owned agents.
	var isAuthenticatedAgent = false
	for _, information := range agents {
		if information.ID == selectedID {
			isAuthenticatedAgent = true
			selectedAgent = information
			break
		}
	}

	// check if selected agent is authorized.
	if !isAuthenticatedAgent {
		// 	If reaches here, it means entered agent id is not authorized.
		return errors.Errorf("%s is not an authorized agent id", selectedID)
	}

	// print selected information
	color.Blue("Selected agent: \n\t %s", selectedAgent.String())

	// check given agent id is ok
	cliConfig.AgentID = selectedID
	agentConfig, err := agent.GetAgentConfig(cliConfig)
	if err != nil {
		return err
	}

	// store selected agent to configuration file.
	viper.Set("AGENT_ID", selectedID)
	viper.Set("DATABASE", agentConfig.Config.Database)
	err = viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to save config. Do you have permissions on the filesystem? ")
	}

	return nil
}

func handleInteractive(agents []agent2.AgentInformation, cliConfig cli.Cli) (
	string, error,
) {
	var (
		err           error
		selectedAgent agent2.AgentInformation
		selectedID    string
	)

	// Ask user how to select agent
	if viper.GetBool("PRETTY_OUTPUT") {
		prompt := promptui.Select{
			Label: "Select an agent",
			Items: []string{"Enter manually", "Choose from list"},
		}

		idx, _, err := prompt.Run()
		if idx == 0 {
			prompt := promptui.Prompt{
				Label:    "Enter Agent ID",
				Validate: validate,
			}

			selectedID, err = prompt.Run()
			if err != nil {
				return "", err
			}

			return selectedID, nil
		}

		// select agent with current cli configuration
		selectedAgent, err = SelectAgentFromList(agents)
		if err != nil {
			return "", err
		}
		selectedID = selectedAgent.ID
	} else {
		reader := bufio.NewReader(os.Stdin)

		color.Blue("Enter agent id: ")
		selectedID, err = reader.ReadString('\n')
		if err != nil {
			return "", errors.WithMessagef(err, "can not read agent id input")
		}

		selectedID = strings.Trim(selectedID, "\n")

		// 	validate input
		err = validate(selectedID)
		if err != nil {
			return "", errors.WithMessagef(err, "can not validate agent id input %s", selectedID)
		}
	}

	return selectedID, nil
}

func validate(input string) error {
	_, err := uuid.Parse(input)
	if err != nil {
		return errors.New("Invalid agent id")
	}
	return nil
}
