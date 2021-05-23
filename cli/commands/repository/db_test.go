package repository

import (
	"testing"
	
	`github.com/spf13/viper`
	"gorm.io/gorm"
	
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	testing database connections
*/

// global constants for file
const ()


func TestConnectToAgentDB(t *testing.T) {
	var (
		cliConfig config.Cli
	)
	
	// get good keys
	viper.SetConfigFile("../../../cli.json")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Error reading config file, %s", err)
	}
	
	err := viper.Unmarshal(&cliConfig)
	
	if err != nil {
		t.Fatalf("unable to decode into map, %v", err)
	}
	
	type args struct {
		dbConfig config.Database
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.DB
		wantErr bool
	}{
		{
			name:    "good credentials",
			args:    args{dbConfig: cliConfig.Database},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := ConnectToAgentDB(tt.args.dbConfig)
				if (err != nil) != tt.wantErr {
					t.Errorf("ConnectToAgentDB() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				_ = got
			},
		)
	}
}
