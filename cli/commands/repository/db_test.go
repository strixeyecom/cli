package repository

import (
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/usestrix/cli/domain/cli"
	"github.com/usestrix/cli/domain/repository"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	testing database connections
*/

// global constants for file
const ()

var (
	disabled = false
)

func TestMain(m *testing.M) {
	// setup test environment
	var (
		err       error
		cliConfig cli.Cli
		dbConfig  repository.Database
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
		log.Fatal(err)
	}

	dbConfig = cliConfig.Database
	dbConfig.SetTestContainerName( "strixeye_test_db")

	// Create a temporary database container for testing
	err = CreateDatabase(cliConfig.Database)
	if err != nil {
		log.Fatal(err)
	}
	// run all tests
	exitCode := m.Run()

	_ = RemoveDatabase(dbConfig)

	// 	exit with test runcode
	os.Exit(exitCode)
}

func TestConnectToAgentDB(t *testing.T) {
	// this test is disabled since integration tests are not my first priority.
	if disabled {
		t.SkipNow()
	}

	var (
		err       error
		cliConfig cli.Cli
	)

	err = viper.Unmarshal(&cliConfig)

	if err != nil {
		t.Fatalf("unable to decode into map, %v", err)
	}

	type args struct {
		dbConfig repository.Database
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

	// Stop created containers

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// disabled because all unit tests are testing this
				// makes no sense until we add more cases to this test.
			},
		)
	}
}
