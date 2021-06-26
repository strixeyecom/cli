package trip

import (
	`fmt`
	`os`
	"testing"
	`time`
	
	`github.com/pkg/errors`
	"github.com/sirupsen/logrus"
	`github.com/spf13/viper`
	`github.com/usestrix/cli/cli/commands/repository`
	repository2 `github.com/usestrix/cli/domain/repository`
	
	`github.com/usestrix/cli/domain/cli`
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
		dbConfig repository2.Database
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
	dbConfig.DBPort = "12348"
	dbConfig.TestContainerName_ = "strixeye_trip_db"
	
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
		err            error
		databaseConfig repository2.Database
		cliConfig      cli.Cli
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

	databaseConfig = cliConfig.Database
	databaseConfig.DBPort = "12348"
	databaseConfig.SetTestContainerName("strixeye_trip_db")
	if err != nil {
		t.Fatalf("unable to decode into map, %v", err)
	}

	type args struct {
		cliConfig repository2.Database
		args      repository2.TripQueryArgs
	}
	tests := []struct {
		name    string
		args    args
		want    []repository2.Trip
		wantErr bool
	}{
		{
			name:    "good credentials",
			args:    args{cliConfig: databaseConfig, args: repository2.TripQueryArgs{Limit: 6}},
			wantErr: false,
		}, {
			name: "verbose output",
			args: args{
				cliConfig: databaseConfig, args: repository2.TripQueryArgs{Limit: 6, Verbose: true},
			},
			wantErr: false,
		}, {
			name: "filter trips by suspect ids",
			args: args{
				cliConfig: databaseConfig,
				args: repository2.TripQueryArgs{
					SuspectIds: []string{
						"3981bb12-8ccc-4493-9884-9d8d46a2ca59", "3981bb12-8ccc-4493-9884-9d8d46a2ca59",
					},
				},
			},
			wantErr: false,
		}, {
			name: "trips newer than T",
			args: args{
				cliConfig: databaseConfig,
				args: repository2.TripQueryArgs{
					SinceTime: int64(time.Now().Second() - 30),
				},
			},
			wantErr: false,
		}, {
			name: "trips to limited endpoints",
			args: args{
				cliConfig: databaseConfig,
				args: repository2.TripQueryArgs{
					Endpoints: []string{"/login"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := get(tt.args.cliConfig, tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				
				_ = got
			},
		)
	}
}
