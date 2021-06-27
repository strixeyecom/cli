package agent

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/usestrix/cli/domain/agent"
	"github.com/usestrix/cli/domain/cli"
)

/*
	Created by aomerk at 6/26/21 for project connector
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func Test_loginToDocker(t *testing.T) {
	// setup test environment
	var (
		err         error
		cliConfig   cli.Cli
		agentConfig agent.AgentInformation
	)

	// get good keys
	viper.SetConfigFile("../../../.env")
	// Try to read from file, but use env variables if non exists. it's fine
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cliConfig)

	agentBytes, err := os.ReadFile("../../../test/agent-info.json")
	if err != nil {
		logrus.Fatal(err)
	}

	err = json.Unmarshal(agentBytes, &agentConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	type args struct {
		agentConfig agent.AgentInformation
		cliConfig   cli.Cli
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pull successfully",
			args: args{
				agentConfig: agentConfig,
				cliConfig:   cliConfig,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := loginToDocker(tt.args.agentConfig, tt.args.cliConfig); (err != nil) != tt.wantErr {
					t.Errorf("loginToDocker() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
