package configure

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/usestrix/cli/api/user/authenticate"
	"github.com/usestrix/cli/domain/config"
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

// NewConfigureUserCommand returns the command for
//  strixeye configure user
// Which lets users authenticate via prompts or flags
func NewConfigureUserCommand() *cobra.Command {
	// serveCmd represents the request command
	var configureUserCommand = &cobra.Command{
		Use:   "user",
		Short: "Change default user to work with",
		Long: `A default user must be selected to work with because StrixEye's CLI is designed to work
with a single StrixEye user's API Key.

Authentication will be executed via this user's API Token' `,
		RunE: selectUser,
	}

	// Add subcommands
	configureUserCommand.AddCommand()

	// Add local flags
	configureUserCommand.Flags().StringP(
		"set-api-key", "s", "",
		`Set api key directly without cli prompt`,
	)

	return configureUserCommand
}

// promptUser prompts api token from user.
func promptUser() (string, error) {
	var (
		err        error
		userAPIKey string
		prompt     *survey.Input
	)

	// get api key from user input
	prompt = &survey.Input{
		Message: "Please enter your user api key\n>",
	}

	// ask user prompt.
	err = survey.AskOne(prompt, &userAPIKey)
	if err != nil {
		return userAPIKey, err
	}

	return userAPIKey, nil
}

// selectUser implements command functionality.
func selectUser(cmd *cobra.Command, args []string) error {
	var (
		err        error
		userAPIKey string
		cliConfig  config.Cli
	)

	// try to get from flags.
	userAPIKey, err = cmd.Flags().GetString("set-api-key")
	if err != nil {
		return err
	}

	// if flags aren't provided, let user choose from a select list.
	if userAPIKey == "" {
		userAPIKey, err = promptUser()
		if err != nil {
			return err
		}
	}

	// try to authenticate user.
	// no need to get the config file, viper supports setting values on the fly.
	viper.Set("USER_API_TOKEN", userAPIKey)
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}

	// ask api for authentication
	err = authenticate.Authenticate(cliConfig)
	if err != nil {
		return errors.Wrap(err, "can not authenticate user with given credential.")
	}

	// store edited viper config to file.
	err = viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to save config. Do you have permissions on the filesystem? ")
	}

	// Write success message.
	color.Blue(
		`Successfully authenticated.
You can start using StrixEye CLI after selecting a default agent if you haven't.
Command for that is:`,
	)
	color.Green("\tstrixeye configure agent")

	return nil
}
