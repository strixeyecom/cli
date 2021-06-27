package ux

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/usestrix/cli/api/user/authenticate"
	"github.com/usestrix/cli/domain/cli"
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

// PromptUserAPIToken prompts api token from user.
func PromptUserAPIToken() (string, error) {
	var (
		err        error
		userAPIKey string
		prompt     *promptui.Prompt
	)

	validate := func(input string) error {
		// TODO add user api token validation logic
		return nil
	}

	if viper.GetBool("PRETTY_OUTPUT") {
		prompt = &promptui.Prompt{
			Label:    "Please enter your User API Token",
			Validate: validate,
		}

		userAPIKey, err = prompt.Run()
		if err != nil {
			return "", err
		}
	} else {
		reader := bufio.NewReader(os.Stdin)

		color.Blue("Enter your StrixEye User API Token: ")
		userAPIKey, err = reader.ReadString('\n')
		if err != nil {
			return "", errors.WithMessagef(err, "can not read user api token input")
		}

		userAPIKey = strings.Trim(userAPIKey, "\n")

		// 	validate input
		err = validate(userAPIKey)
		if err != nil {
			return "", errors.WithMessagef(err, "can not validate user api token input %s", userAPIKey)
		}
	}

	return userAPIKey, nil
}

// SetupUser implements command functionality.
func SetupUser(cmd *cobra.Command, args []string) error {
	var (
		err        error
		userAPIKey string
		cliConfig  cli.Cli
	)

	// try to get from flags.
	userAPIKey, err = cmd.Flags().GetString("user-api-token")
	if err != nil {
		return err
	}

	// if flags aren't provided, let user choose from a select list.
	if userAPIKey == "" {
		userAPIKey, err = PromptUserAPIToken()
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
