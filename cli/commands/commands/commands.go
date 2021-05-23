package commands

import (
	"fmt"
	"strings"
	
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	
	"github.com/usestrix/cli/cli/commands/configure"
	"github.com/usestrix/cli/cli/commands/trip"
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	Both config initialization and cli initialization takes place in here.
*/

// global constants for file
const (
	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = "cli"
	
	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to STING_NUMBER.
	envPrefix = "STRIXEYE"
)

// global variables (not cool) for this file
var (
	cfgFile string
)

// NewStrixeyeCommand is the highest command in the hierarchy and all commands root from it.
//nolint:funlen
func NewStrixeyeCommand() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	// Define our command
	rootCmd := &cobra.Command{
		Use:   "strixeye",
		Short: "The StrixEye Command Line Interface",
		Long:  `Inspect and Manage your agents with strixeye cli from anywhere.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var (
				cliConfig config.Cli
				err       error
			)
			
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			err = handleConfig(cmd)
			if err != nil {
				return err
			}
			
			// get values from viper
			a := viper.GetViper()
			_ = a
			
			// unmarshal into config object
			err = viper.Unmarshal(&cliConfig)
			if err != nil {
				return err
			}
			// if it is not a valid config file, force user to configure strixeye cli first.
			// err = cliConfig.Validate()
			// if err != nil {
			// 	return fmt.Errorf(
			// 		"validating configuration: %w\nplease check your configuration file or run"+
			// 			"\n\t>  strixeye configure", err,
			// 	)
			// }
			
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			// out := cmd.OutOrStdout()
		},
	}
	
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.strixeye/."+
			"config.yaml)",
	)
	
	// Add subcommands
	rootCmd.AddCommand(
		trip.NewTripCommand(),
		configure.NewConfigureCommand(),
	)
	
	// Add flags
	
	return rootCmd
}

// initializeConfig tries to open config and validate it. If a bad config is given,
// tells user possible solutions.
func handleConfig(cmd *cobra.Command) error {
	err := initializeConfig(cmd)
	if err == nil {
		return nil
	}
	
	// 	cli config not found or no permission to read.
	color.Red("Please create a config file.")
	return nil
}

func initializeConfig(cmd *cobra.Command) error {
	// Set the base name of the config file, without the file extension.
	viper.SetConfigName(defaultConfigFilename)
	
	// viper.SetDefault("API_URL", "api.strixeye.com")
	
	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		
		// Search config in home directory with name ".cli" (without extension).
		
		viper.AddConfigPath(home + "/.strixeye")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/strixeye")
		
	}
	
	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		var _t0 viper.ConfigFileNotFoundError
		if ok := errors.Is(err, _t0); !ok {
			return err
		}
	}
	
	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(envPrefix)
	
	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	viper.AutomaticEnv()
	
	// Bind the current command's flags to viper
	bindFlags(cmd, viper.GetViper())
	
	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(
		func(f *pflag.Flag) {
			// Environment variables can't have dashes in them, so bind them to their equivalent
			// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
			if strings.Contains(f.Name, "-") {
				envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
				err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
				cobra.CheckErr(err)
			}
			
			// Apply the viper config value to the flag when the flag is not set and viper has a value
			if !f.Changed && v.IsSet(f.Name) {
				val := v.Get(f.Name)
				err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
				cobra.CheckErr(err)
			}
		},
	)
}
