package agent

import (
	"fmt"
	"testing"
	
	"github.com/spf13/viper"
	
	`github.com/usestrix/cli/domain/agent`
	`github.com/usestrix/cli/domain/cli`
)

/*
	Created by aomerk at 5/22/21 for project cli
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func TestGetAgents(t *testing.T) {
	var (
		err       error
		cliConfig cli.Cli
	)
	viper.SetConfigFile("../../../.env")
	// Try to read from file, but use env variables if non exists. it's fine
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}
	viper.AutomaticEnv()
	
	err = viper.Unmarshal(&cliConfig)
	
	if err != nil {
		fmt.Printf("Unable to decode into map, %v", err)
	}
	type args struct {
		cliConfig cli.Cli
	}
	tests := []struct {
		name    string
		args    args
		want    []agent.AgentInformation
		wantErr bool
	}{
		{
			name: "Get all agents",
			args: args{
				cliConfig: cliConfig,
			},
			wantErr: false,
		}, {
			name: "Authentication failure",
			args: args{
				cliConfig: cli.Cli{UserAPIToken: "fake-token", APIUrl: cliConfig.APIUrl},
			},
			wantErr: true,
		}, {
			name: "Bad Request",
			args: args{
				cliConfig: cli.Cli{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := GetAgents(tt.args.cliConfig)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAgents() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				_ = got
			},
		)
	}
}
