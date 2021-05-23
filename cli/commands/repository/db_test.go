package repository

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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

// global variables (not cool) for this file
var (
	dbConfig_ config.Database
)

func TestMain(m *testing.M) {
	var (
		err error
	)

	// initialize test environment
	err = godotenv.Load("cli.json")
	if err != nil {
		logrus.Fatal(err)
	}

	// create a real database instance
	dbConfig_ = config.Database{
		DBName: os.Getenv("DB_NAME"),
		DBPort: os.Getenv("DB_PORT"),
		DBAddr: os.Getenv("DB_ADDR"),
		DBPass: os.Getenv("DB_PASS"),
		DBUser: os.Getenv("DB_USER"),
	}

	err = dbConfig_.Validate()
	if err != nil {
		logrus.Fatal(err)
	}

	// 	run tests
	os.Exit(m.Run())
}
func TestConnectToAgentDB(t *testing.T) {
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
			args:    args{dbConfig: dbConfig_},
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
