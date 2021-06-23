package ux

import (
	`github.com/fatih/color`
	`github.com/google/uuid`
	"github.com/manifoldco/promptui"
	`github.com/pkg/errors`
	`github.com/spf13/cobra`
	`github.com/spf13/viper`
	`github.com/usestrix/cli/api/user/agent`
	agent2 `github.com/usestrix/cli/domain/agent`
	`github.com/usestrix/cli/domain/cli`
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
	
	selectedID = viper.GetString("CURRENT_AGENT_ID")
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
	color.Blue("Selected agent: %s", selectedAgent.String())
	
	// check given agent id is ok
	cliConfig.CurrentAgentID = selectedID
	_, err = agent.GetAgentConfig(cliConfig)
	if err != nil {
		return err
	}
	
	// store selected agent to configuration file.
	viper.Set("CURRENT_AGENT_ID", selectedID)
	err = viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to save config. Do you have permissions on the filesystem? ")
	}
	
	return nil
}

func handleInteractive(agents []agent2.AgentInformation, cliConfig cli.Cli,) (
	string, error,
) {
	var (
		err        error
		selectedID string
	)
	
	prompt := promptui.Select{
		Label: "Select an agent",
		Items: []string{"Enter manually", "Choose from list"},
	}
	
	idx, _, err := prompt.Run()
	if idx == 0 {
		validate := func(input string) error {
			_, err = uuid.Parse(input)
			if err != nil {
				return errors.New("Invalid agent id")
			}
			return nil
		}
		
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
	selectedAgent, err := SelectAgentFromList(agents)
	if err != nil {
		return "", err
	}
	
	selectedID = selectedAgent.ID
	
	return selectedID, nil
}
