package suspect

import (
	"testing"
	
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	
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
		dbConfig  config.Database
		cliConfig config.Cli
	)
	// get good keys
	
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
	
	dbConfig = cliConfig.Database
	

	err = dbConfig.Validate()
	if err != nil {
		logrus.Fatal(err)
	}
	
	type args struct {
		dbConfig config.Database
		args     QueryArgs
	}
	tests := []struct {
		name    string
		args    args
		want    []Suspect
		wantErr bool
	}{
		{
			name:    "good credentials",
			args:    args{dbConfig: dbConfig, args: QueryArgs{Limit: 6}},
			wantErr: false,
		}, {
			name: "suspects with score bigger than",
			args: args{
				dbConfig: dbConfig,
				args: QueryArgs{
					Score: 5,
				},
			},
			wantErr: false,
		}, {
			name: "filter by suspect ids",
			args: args{
				dbConfig: dbConfig,
				args: QueryArgs{
					SuspectIds: []string{
						"3981bb12-8ccc-4493-9884-9d8d46a2ca59", "3981bb12-8ccc-4493-9884-9d8d46a2ca59",
					},
				},
			},
			wantErr: false,
		}, {
			name: "suspect newer than T",
			args: args{
				dbConfig: dbConfig,
				args: QueryArgs{
					SinceTime: 1621508903188,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := Get(tt.args.dbConfig, tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				
				_ = got
			},
		)
	}
}
