package trip

import (
	"testing"
	
	"github.com/sirupsen/logrus"
	`github.com/spf13/viper`
	
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
 
 */

// global constants for file
const ()


func TestGet(t *testing.T) {
	var (
		err       error
		cliConfig config.Cli
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
	
	if err != nil {
		t.Fatalf("unable to decode into map, %v", err)
	}
	
	type args struct {
		cliConfig config.Cli
		args      QueryArgs
	}
	tests := []struct {
		name    string
		args    args
		want    []Trip
		wantErr bool
	}{
		{
			name:    "good credentials",
			args:    args{cliConfig: cliConfig, args: QueryArgs{Limit: 6}},
			wantErr: false,
		},{
			name:    "verbose output",
			args:    args{cliConfig: cliConfig, args: QueryArgs{Limit: 6,Verbose: true}},
			wantErr: false,
		}, {
			name: "filter trips by suspect ids",
			args: args{
				cliConfig: cliConfig,
				args: QueryArgs{
					SuspectIds: []string{
						"3981bb12-8ccc-4493-9884-9d8d46a2ca59", "3981bb12-8ccc-4493-9884-9d8d46a2ca59",
					},
				},
			},
			wantErr: false,
		}, {
			name: "trips newer than T",
			args: args{
				cliConfig: cliConfig,
				args: QueryArgs{
					SinceTime: 1621508903188,
				},
			},
			wantErr: false,
		}, {
			name: "trips to limited endpoints",
			args: args{
				cliConfig: cliConfig,
				args: QueryArgs{
					Endpoints: []string{"/login"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := Get(tt.args.cliConfig, tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				
				_ = got
			},
		)
	}
}
