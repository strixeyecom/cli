package commands

import (
	`encoding/json`
	`fmt`
	
	`github.com/fatih/color`
	`github.com/go-yaml/yaml`
	`github.com/pelletier/go-toml`
	`github.com/pkg/errors`
	`github.com/spf13/cobra`
	`github.com/spf13/viper`
	`github.com/usestrix/cli/domain/cli`
)

/*
	Created by aomerk at 8/3/21 for project strixeye
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// InspectCommand extracts information from CLI and Agent configs and redirects to output device.
func InspectCommand() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "inspect",
		Short: "Command to show agent information",
		Long: `Command to show agent information

This command will show configuration information your current agent on this machine.

strixeye agent inspect
`,
		RunE: inspectCmd,
	}
	
	// set up flags
	checkCmd.Flags().String(
		"format", "json", "--format if you want to output information in a specified format like json, "+
			"yaml or toml",
	)
	return checkCmd
}

func inspectCmd(cmd *cobra.Command, args []string) error {
	var (
		cliConfig cli.Cli
		err       error
	)
	
	// get cli config
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}
	
	// show only requested fields
	if len(args) > 0 {
		// iterate all wanted fields
		for _, arg := range args {
			str := viper.GetString(arg)
			color.Blue("%s", str)
		}
		return nil
	}
	
	// output in requested format
	fmtFlag, err := cmd.Flags().GetString("format")
	if err != nil {
		return errors.Wrap(err, "can not inspect, bad format value")
	}
	conf, err := marshalToFormat(cliConfig, fmtFlag)
	if err != nil {
		return errors.Wrap(err, "can not marshal cli config while inspecting")
	}
	fmt.Println(string(conf))
	return nil
}

func marshalToFormat(cliConfig cli.Cli, fmt string) ([]byte, error) {
	switch fmt {
	case "json":
		return json.MarshalIndent(cliConfig, "", "\t")
	case "yaml":
		return yaml.Marshal(cliConfig)
	case "toml":
		return toml.Marshal(cliConfig)
	}
	return nil, errors.New("unsupported format")
}
