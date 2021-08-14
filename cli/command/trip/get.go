package trip

import (
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	
	userconfig "github.com/strixeyecom/cli/api/user/agent"
	`github.com/strixeyecom/cli/cli/command/repository`
	`github.com/strixeyecom/cli/domain/cli`
	models `github.com/strixeyecom/cli/domain/repository`
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	get request/response pairs from your agent database
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// GetCommand returns cli command to query trips.
// Command to query trips on your agent
// StrixEye agent logs requests/responses of your visitors for a span of time.
// If you are our customer, then you know that StrixEye agent runs on your internal network and doesn't leak
// any data outside of your network/its network because of privacy and security concerns.
//
// With this subcommand, you can inspect your logs get is a subcommand of trip command where you can query
// trips on your agent, without leaking any sensitive data outside of your network.
func GetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Command to query trips on your agent",
		Long: `Command to query trips on your agent
StrixEye agent logs requests/responses of your visitors for a span of time.
If you are our customer, then you know that StrixEye agent runs on your internal network and doesn't leak
any data outside of your network/its network because of privacy and security concerns.

With this subcommand, you can inspect your logs get is a subcommand of trip command where you can query
trips on your agent, without leaking any sensitive data outside of your network`,
		RunE: getTripCmd,
	}
	
	// declaring local flags used by get trip commands.
	getCmd.Flags().StringSliceP(
		"suspects", "s", nil, "Comma separated values of suspect uuids. "+
			"--suspects suspect-id,suspect-id2",
	)
	
	// declaring local flags used by get trip commands.
	getCmd.Flags().StringSliceP(
		"ids", "t", nil, "Comma separated values of trip uuids. "+
			"--ids trip-id,trip-id2",
	)
	getCmd.Flags().StringSliceP(
		"endpoints", "e", nil,
		"Endpoints that you want to display. Comma separated list of endpoints. "+
			"--endpoints /login,/logout",
	)
	getCmd.Flags().IntP("limit", "l", 5, "Max number of trips you want to be displayed --limit 5")
	
	getCmd.Flags().BoolP(
		"verbose", "v", true, "To hide field values like headers, say --verbose=false",
	)
	getCmd.Flags().IntP(
		"since", "i", 0,
		"Queries only suspicions after given time --since [epoch in seconds]  You can get current"+
			" timestamp with date +%s",
	)
	return getCmd
}

// getTripCmd implements GetCommand logic.
func getTripCmd(cmd *cobra.Command, _ []string) error {
	var (
		cliConfig cli.Cli
		err       error
	)
	
	// get cli config for authentication
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}
	
	queryArgs := models.TripQueryArgs{Limit: 1}
	
	// parse and set list of suspects to be queried
	suspects, err := cmd.Flags().GetStringSlice("suspects")
	if err != nil {
		return err
	} else if len(suspects) > 0 {
		queryArgs.SuspectIds = suspects
	}
	
	// parse and set list of endpoints to be queried
	endpoints, err := cmd.Flags().GetStringSlice("endpoints")
	if err != nil {
		return err
	} else if len(endpoints) > 0 {
		queryArgs.Endpoints = endpoints
	}
	
	// parse and set list of suspicions to be queried
	tripIds, err := cmd.Flags().GetStringSlice("ids")
	if err != nil {
		return err
	} else if len(tripIds) > 0 {
		queryArgs.TripsIds = tripIds
	}
	
	// parse max limit of rows displayed.
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return err
		
	} else if limit > 0 {
		queryArgs.Limit = limit
	}
	
	// parse verboseness flag.
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	} else if limit > 0 {
		queryArgs.Verbose = verbose
	}

	// parse oldest timestamp to be queries
	sinceTime, err := cmd.Flags().GetInt("since")
	if err != nil {
		return err

	} else if sinceTime > 0 {
		queryArgs.SinceTime = int64(sinceTime)
	}

	// get trips with parsed query arguments for this subcommand
	trips, err := Get(cliConfig, queryArgs)
	if err != nil {
		return err
	}

	// print out query settings
	color.Blue(queryArgs.String())

	if trips == nil || len(trips) == 0 {
		color.Blue("0 result.")
		return nil
	}

	// print out result
	_, err = pp.Print(trips)
	if err != nil {
		return err
	}

	return nil
}

// Get is a temporary method to satisfy the authentication process.
func Get(cliConfig cli.Cli, args models.TripQueryArgs) ([]models.Trip, error) {
	var (
		dbConfig models.Database
	)
	
	// If user wants to override db config with local information, use that.
	if cliConfig.Database.Validate() == nil && cliConfig.Database.OverrideRemoteConfig {
		color.Blue("Using local database config.")
		dbConfig = cliConfig.Database
	} else {
		agentConfig, err := userconfig.GetAgentConfig(cliConfig)
		if err != nil {
			return nil, err
		}
		dbConfig = agentConfig.Config.Database
	}
	
	return get(dbConfig, args)
}

// Get retrieves all trips that matches given query args. Check out trips.
// TripQueryArgs for more information about existing filters.
func get(dbConfig models.Database, args models.TripQueryArgs) ([]models.Trip, error) {
	var (
		err    error
		db     *gorm.DB
		result []models.Trip
	)
	// connect to database
	db, err = repository.ConnectToAgentDB(dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "can not establish connection to agent database")
	}
	
	// nobody wants to retrieve all hundreds of thousands of results.
	if args.Limit == 0 {
		args.Limit = 10
	}
	tx := db.Limit(args.Limit)
	
	if args.Verbose {
		tx = tx.Preload("Request.Header")
	}
	// filter by endpoints
	if args.Endpoints != nil {
		tx = tx.Joins("Request").Where("raw_uri IN ? ", args.Endpoints)
	} else {
		// preload only first level for now
		tx = tx.Preload(clause.Associations)
	}
	
	// filter by suspect ids
	if args.SuspectIds != nil {
		tx = tx.Where("profile_id IN ?", args.SuspectIds)
	}
	
	// filter by suspect ids
	if args.TripsIds != nil {
		tx = tx.Where(args.TripsIds)
	}
	
	// filter by created after
	if args.SinceTime != 0 {
		// convert seconds to milliseconds
		milliseconds := args.SinceTime * 1e3
		tx = tx.Where("created_at > ?", milliseconds)
	}
	
	// find all suspects that matches args criteria
	tx = tx.Find(&result)
	err = tx.Error
	return result, err
}
