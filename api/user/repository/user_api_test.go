package repository

import (
	"io"
	"net/http"
	"testing"
	
	"github.com/spf13/viper"
	
	`github.com/usestrix/cli/domain/cli`
)

/*
	Created by aomerk at 5/22/21 for project cli
*/

/*
	test for user api communication
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func TestUserAPIRequest(t *testing.T) {
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
		t.Errorf("unable to decode into map, %v", err)
	}

	type args struct {
		method   string
		endpoint string
		body     io.Reader
		apiToken string
		apiURL   string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "Get with no body",
			args: args{
				apiToken: cliConfig.UserAPIToken,
				apiURL:   cliConfig.APIDomain,
				body:     nil,
				endpoint: "/agents",
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
				got, err := UserAPIRequest(
					tt.args.method, tt.args.endpoint, tt.args.body, tt.args.apiToken, tt.args.apiURL,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("UserAPIRequest() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				_ = got
			},
		)
	}
}
