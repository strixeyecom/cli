package suspicion

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
		want    []Suspicion
		wantErr bool
	}{
		{
			name:    "good credentials",
			args:    args{cliConfig: cliConfig, args: QueryArgs{Limit: 6}},
			wantErr: false,
		}, {
			name: "limit suspicion ids",
			args: args{
				cliConfig: cliConfig,
				args: QueryArgs{
					SuspicionIds: []string{
						"15423ea9-d658-4cc2-ba21-37acd0f90c07", "40d91d87-3103-4ba9-9aae-55e3b05b1d70",
					},
				},
			},
			wantErr: false,
		}, {
			name: "filter by suspect ids",
			args: args{
				cliConfig: cliConfig,
				args: QueryArgs{
					SuspectIds: []string{
						"3981bb12-8ccc-4493-9884-9d8d46a2ca59", "3981bb12-8ccc-4493-9884-9d8d46a2ca59",
					},
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
