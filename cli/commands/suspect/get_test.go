package suspect

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/strixeyecom/cli/cli/commands/repository"
	models "github.com/strixeyecom/cli/domain/repository"

	"github.com/strixeyecom/cli/domain/cli"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*

 */

// global constants for file
const ()

func TestMain(m *testing.M) {
	exitCode, err := wrapper(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitCode)
}

func wrapper(m *testing.M) (int, error) {
	var (
		exitCode = 1
		dbConfig models.Database
	)

	defer func() {
		// Clean up on exit.
		_ = repository.RemoveDatabase(dbConfig)
	}()

	// setup test environment
	var (
		err       error
		cliConfig cli.Cli
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
		return 1, err
	}

	dbConfig = cliConfig.Database
	dbConfig.DBPort = "12346"
	dbConfig.SetTestContainerName("strixeye_suspect_db")

	// delete all database data just in case
	_ = repository.RemoveDatabase(dbConfig)

	// Create a temporary database container for testing
	err = repository.CreateDatabaseIFNotExists(dbConfig)
	if err != nil {
		return 1, err
	}

	// Setup temporary database
	err = repository.SetupDatabase(dbConfig)
	if err != nil {
		return 1, err
	}

	// run all tests
	exitCode = m.Run()
	if exitCode != 0 {
		return exitCode, errors.New("tests are failed")
	}

	return exitCode, nil
}

func TestGet(t *testing.T) {
	var (
		err       error
		dbConfig  models.Database
		cliConfig cli.Cli
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
	dbConfig.DBPort = "12346"
	dbConfig.SetTestContainerName("strixeye_suspect_db")

	err = dbConfig.Validate()
	if err != nil {
		logrus.Fatal(err)
	}

	type args struct {
		dbConfig models.Database
		args     models.SuspectQueryArgs
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Suspect
		wantErr bool
	}{
		{
			name:    "good credentials",
			args:    args{dbConfig: dbConfig, args: models.SuspectQueryArgs{Limit: 6}},
			wantErr: false,
		}, {
			name: "suspects with score bigger than",
			args: args{
				dbConfig: dbConfig,
				args: models.SuspectQueryArgs{
					MinScore: 5,
				},
			},
			wantErr: false,
		}, {
			name: "filter by suspect ids",
			args: args{
				dbConfig: dbConfig,
				args: models.SuspectQueryArgs{
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
				args: models.SuspectQueryArgs{
					SinceTime: time.Now().UnixNano() - 5*time.Second.Nanoseconds(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := get(tt.args.dbConfig, tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				_ = got
			},
		)
	}
}
