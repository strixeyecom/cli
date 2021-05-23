package agent

import (
	"testing"
	
	"github.com/spf13/viper"
	
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

// global constants for file
const ()

func TestGetAgent(t *testing.T) {
	var (
		// get good keys
		err error
		cliConfig config.Cli
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
		t.Fatalf("unable to decode into map, %v", err)
	}
	
	type args struct {
		userAPIToken string
		agentID      string
		apiURL       string
	}
	tests := []struct {
		name    string
		args    args
		want    config.AgentInformation
		wantErr bool
	}{
		{
			name: "get with good token",
			args: args{
				userAPIToken: cliConfig.UserAPIToken,
				agentID:      cliConfig.CurrentAgentID,
				apiURL:       cliConfig.APIUrl,
			},
			wantErr: false,
		}, {
			name: "get with good agent with bad token",
			args: args{
				userAPIToken: "fake-token",
				agentID:      cliConfig.CurrentAgentID,
				apiURL:       cliConfig.APIUrl,
			},
			wantErr: true,
		},
		{
			name: "get bad agent with good token",
			args: args{
				userAPIToken: cliConfig.UserAPIToken,
				agentID:      "fake-agent",
				apiURL:       cliConfig.APIUrl,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getAgent(tt.args.userAPIToken, tt.args.apiURL, tt.args.agentID)
				if (err != nil) != tt.wantErr {
					t.Errorf("FetchAgentById() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err == nil {
					_ = got.String()
					_ = got
				}
			},
		)
	}
}
