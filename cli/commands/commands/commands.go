package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	agent2 "github.com/usestrix/cli/domain/agent"
	"github.com/usestrix/cli/domain/consts"
	
	"github.com/usestrix/cli/cli/commands/agent"
	"github.com/usestrix/cli/cli/commands/configure"
	"github.com/usestrix/cli/cli/commands/suspect"
	"github.com/usestrix/cli/cli/commands/suspicion"
	"github.com/usestrix/cli/cli/commands/trip"
	"github.com/usestrix/cli/domain/cli"
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
	defaultConfigFileType = "yaml"
	
	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to STING_NUMBER.
	envPrefix = "STRIXEYE"
)

// global variables (not cool) for this file
var (
	cfgFile        string
	apiDomain      string
	downloadDomain string
	userAPIToken   string
	agentID        string
)

// NewStrixeyeCommand is the highest command in the hierarchy and all commands root from it.
//nolint:funlen
func NewStrixeyeCommand() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	// Define our command
	rootCmd := &cobra.Command{
		Use:   "strixeye",
		Short: "The StrixEye Command Line Interface",
		Long:  `Inspect and Manage your agents with StrixEye CLI from anywhere.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var (
				cliConfig cli.Cli
				err       error
			)
			
			// don't allow non strixeye or non-root users to run StrixEye CLI
			// checkUser(cmd)
			
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			err = handleConfig(cmd)
			if err != nil {
				return errors.WithMessage(err, "can not read config")
			}
			
			// get values from viper
			a := viper.GetViper()
			_ = a
			
			// unmarshal into config object
			err = viper.Unmarshal(&cliConfig)
			if err != nil {
				return err
			}
			
			return nil
		},
		RunE: ShowHelp(os.Stdout),
	}
	
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.strixeye/."+
			defaultConfigFilename+"."+defaultConfigFileType+")",
	)
	rootCmd.PersistentFlags().StringVar(
		&apiDomain, "api-domain", "", "api.strixeye.com",
	)
	
	rootCmd.PersistentFlags().StringVar(
		&downloadDomain, "download-domain", "", "downloads.strixeye.com",
	)
	rootCmd.PersistentFlags().StringVar(
		&userAPIToken, "user-api-token", "", "",
	)
	
	rootCmd.PersistentFlags().StringVar(
		&agentID, "agent-id", "", "",
	)
	
	// Add subcommands
	rootCmd.AddCommand(
		trip.NewTripCommand(),
		configure.NewConfigureCommand(),
		suspicion.NewSuspicionCommand(),
		suspect.NewSuspectCommand(),
		agent.NewAgentCommand(),
	)
	
	// Add flags
	
	return rootCmd
}

func checkUser(cmd *cobra.Command) {
	// if this is the first time, it needs to be the root user.
	// Because strixeye install command runs as a root and if you keep strixeye cli config in a
	// user owned directory, like any directory under $HOME,
	// it won't be accessible by root user by $HOME path, because for example in Linux,
	// $HOME for root user is /root, and that means the config file is under /root/strixeye-cli.
	//
	// Because of this, we need to put it in a non-user based directory. However,
	// it is totally fine to own the directory.
	if !agent2.IsRootUser() {
		color.Red(
			`For now, StrixEye doesn't support non-root users for your security.`,
		)
		os.Exit(1)
	}
}

// initializeConfig tries to open config and validate it. If a bad config is given,
// tells user possible solutions.
func handleConfig(cmd *cobra.Command) error {
	err := initializeConfig(cmd)
	if err == nil {
		return nil
	}
	if cmd.Parent() != nil {
		if cmd.Parent().Use != "configure" {
			// 	cli config not found or no permission to read.
			return err
		}
	} else {
		// 	cli config not found or no permission to read.
		return err
	}
	
	return nil
}

//nolint:funlen
func initializeConfig(cmd *cobra.Command) error {
	var (
		err error
	)
	
	// Set the base name of the config file, without the file extension.
	viper.SetConfigName(defaultConfigFilename)
	viper.SetConfigType(defaultConfigFileType)
	
	setDefaultConfig()
	
	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := getHome()
		cobra.CheckErr(err)
		
		// Search config in home directory with name ".cli" (without extension).
		configPath := home + "/.strixeye"
		viper.AddConfigPath(configPath)
		// viper.AddConfigPath(".")
		
		cfgFile = filepath.Join(configPath, defaultConfigFilename, defaultConfigFileType)
		
		// create default config directory since we are going to use this anyway.
		_, statErr := os.Stat(configPath)
		
		if os.IsNotExist(statErr) {
			// Than, create the directory with root perms only. Actually,
			// a permission like 0666 would prevent the user from `chown`ing the directory,
			// but this decision is not up to me.
			err = os.Mkdir(configPath, 0600)
			if err != nil {
				if os.IsPermission(err) {
					return fmt.Errorf(
						`please set permissions of the directory, eg. using
$ chown -R $USER %s   `, consts.CLIConfigDir,
					)
				}
				return err
			}
			
			// If you are here, then the configuration directories are created and owned by the current process owner
			color.Blue(
				"Successfully set up strixeye cli. "+
					"Please own %s if you want to use it as a non-root user", consts.CLIConfigDir,
			)
		} else if os.IsPermission(err) {
			return errors.WithMessagef(
				err, `Please set permissions of the directory, eg. using
$ chown -R $USER %s`, consts.CLIConfigDir,
			)
		}
	}
	
	// after creating file, we can start using default config file.
	
	// and store it as default
	err = viper.SafeWriteConfig()
	if err != nil {
		// this is not my fault. It is either poor documentation or it is my fault.
		if !(strings.Contains(err.Error(), "Config File") && strings.Contains(
			err.Error(),
			"Already Exists",
		)) {
			// handle failed write
			return err
		}
	}
	// }
	
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

func getHome() (string, error) {
	// 	get os independent home directory
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	
	// 	handle linux sudo fallback
	if home == "/root" {
		sudoHome := os.Getenv("SUDO_HOME")
		if sudoHome != "" {
			color.Blue("Process is running via sudo")
			home = sudoHome
		}
	}
	
	return home, nil
}
func setDefaultConfig() {
	viper.SetDefault("API_DOMAIN", consts.APIHost)
	viper.SetDefault("DOWNLOAD_DOMAIN", consts.DownloadHost)
	viper.SetDefault("DOCKER_REGISTRY", consts.DockerRegistry)
	
	viper.SetDefault("AGENT_ID", "")
	viper.SetDefault("USER_API_TOKEN", "")
	
	viper.SetDefault("PRETTY_OUTPUT", false)
	
	viper.SetDefault("DATABASE.DB_ADDR", "127.0.0.1")
	viper.SetDefault("DATABASE.DB_USER", "")
	viper.SetDefault("DATABASE.DB_PASS", "")
	viper.SetDefault("DATABASE.DB_NAME", "strixeye")
	viper.SetDefault("DATABASE.DB_PORT", "")
	viper.SetDefault("DATABASE.DB_OVERRIDE", true)
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(
		func(f *pflag.Flag) {
			// Environment variables can't have dashes in them, so bind them to their equivalent
			// keys with underscores, e.g. --favorite-color to STRIXEYE_FAVORITE_COLOR
			if strings.Contains(f.Name, "-") {
				envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
				err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
				cobra.CheckErr(err)
			}
			
			if f.Changed {
				envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
				v.Set(envVarSuffix, f.Value)
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

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOut(err)
		d := color.New(color.FgBlue, color.Bold)
		_, _ = d.Print(strixeyeAscii)
		d = color.New(color.FgGreen, color.Bold)
		_, _ = d.Print(normalOwl)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}

const (
	normalOwl = `
  ___
 (o,o)
 {'"'}
 -"-"-
`
	strixeyeAscii = `
  ___   _           _         ___
 / __| | |_   _ _  (_) __ __ | __|  _  _   ___
 \__ \ |  _| | '_| | | \ \ / | _|  | || | / -_)
 |___/  \__| |_|   |_| /_\_\ |___|  \_, | \___|
                                    |__/
`
)
