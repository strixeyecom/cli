package trip

import (
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	userconfig "github.com/usestrix/cli/api/user/config"
	"github.com/usestrix/cli/cli/commands/repository"
	"github.com/usestrix/cli/domain/config"
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
	getCmd.Flags().StringSliceP(
		"endpoints", "e", nil,
		"Endpoints that you want to display. Comma separated list of endpoints. "+
			"--endpoints /login,/logout",
	)
	getCmd.Flags().IntP("limit", "l", 5, "Max number of trips you want to be displayed --limit 5")

	return getCmd
}

// getTripCmd implements GetCommand logic.
func getTripCmd(cmd *cobra.Command, _ []string) error {
	var (
		cliConfig config.Cli
		err       error
	)

	// get cli config for authentication
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}

	queryArgs := QueryArgs{Limit: 1}

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

	// parse max limit of rows displayed.
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return err

	} else if limit > 0 {
		queryArgs.Limit = limit
	}

	// get trips with parsed query arguments for this subcommand
	trips, err := Get(cliConfig, queryArgs)
	if err != nil {
		return err

	}

	// marshal result with colors
	data, err := prettyjson.Marshal(trips)
	if err != nil {
		return err

	}

	// print out query settings
	color.Blue(queryArgs.String())

	// print out result
	color.Blue("%s", string(data))

	return nil
}

// Get is a temporary method to satisfy the authentication process.
func Get(cliConfig config.Cli, args QueryArgs) ([]Trip, error) {

	agentConfig, err := userconfig.GetAgentConfig(cliConfig)
	if err != nil {
		return nil, err
	}

	return get(agentConfig.Config.Database, args)
}

// Get retrieves all trips that matches given query args. Check out trips.
// QueryArgs for more information about existing filters.
func get(dbConfig config.Database, args QueryArgs) ([]Trip, error) {
	var (
		err    error
		db     *gorm.DB
		result []Trip
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
		tx = tx.Where("created_at > ?", args.SinceTime)
	}

	// find all suspects that matches args criteria
	tx = tx.Find(&result)
	err = tx.Error
	return result, err
}
