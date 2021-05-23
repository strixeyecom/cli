package user

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/usestrix/cli/domain/config"
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
	var cliConfig config.Cli
	viper.SetConfigFile("cli.json")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&cliConfig)

	if err != nil {
		fmt.Printf("Unable to decode into map, %v", err)
	}
	type args struct {
		apiToken string
		apiURL   string
	}
	tests := []struct {
		name    string
		args    args
		want    []config.AgentInformation
		wantErr bool
	}{
		{
			name: "Get all agents",
			args: args{
				apiToken: cliConfig.UserAPIToken,
				apiURL:   cliConfig.APIUrl,
			},
			wantErr: false,
		}, {
			name: "Authentication failure",
			args: args{
				apiToken: "fake-token",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getAgents(tt.args.apiToken, tt.args.apiURL)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAgents() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				_ = got
			},
		)
	}
}
