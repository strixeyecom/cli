package authenticate

import (
	"testing"
	
	"github.com/spf13/viper"
	
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	User API authentication tests.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func TestAuthenticate(t *testing.T) {
	var (
		err  error
		conf config.Cli
	)
	
	viper.SetConfigFile("../../../.env")
	// Try to read from file, but use env variables if non exists. it's fine
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}
	viper.AutomaticEnv()
	
	err = viper.Unmarshal(&conf)
	if err != nil {
		t.Fatal(err)
	}
	
	type args struct {
		cliConfig config.Cli
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Successfully authenticate",
			args:    args{cliConfig: conf},
			wantErr: false,
		},
		{
			name:    "Bad authenticate",
			args:    args{cliConfig: config.Cli{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := Authenticate(tt.args.cliConfig); (err != nil) != tt.wantErr {
					t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
