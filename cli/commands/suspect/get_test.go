package suspect

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/usestrix/cli/cli/commands/repository"
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*

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
	err = godotenv.Load(".env")
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

func TestSuspectStruct(t *testing.T) {
	db, err := repository.ConnectToAgentDB(dbConfig_)
	if err != nil {
		t.Error(err)
	}
	_ = db
	// be sure that these structs play well with
}

func TestGet(t *testing.T) {
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
			args:    args{dbConfig: dbConfig_, args: QueryArgs{Limit: 6}},
			wantErr: false,
		}, {
			name: "suspects with score bigger than",
			args: args{
				dbConfig: dbConfig_,
				args: QueryArgs{
					Score: 5,
				},
			},
			wantErr: false,
		}, {
			name: "filter by suspect ids",
			args: args{
				dbConfig: dbConfig_,
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
				dbConfig: dbConfig_,
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
