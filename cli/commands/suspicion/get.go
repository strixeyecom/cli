package suspicion

import (
	`github.com/fatih/color`
	`github.com/hokaccha/go-prettyjson`
	"github.com/pkg/errors"
	`github.com/spf13/cobra`
	`github.com/spf13/viper`
	"gorm.io/gorm"
	
	userconfig `github.com/usestrix/cli/api/user/agent`
	"github.com/usestrix/cli/cli/commands/repository"
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	get suspicions from your agent database
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// GetCommand returns cli command to query suspicions.
// Command to query suspicions on your agent
// StrixEye agent logs suspicions(anomalies) of your visitors for a span of time.
// If you are our customer, then you know that StrixEye agent runs on your internal network and doesn't leak
// any data outside of your network/its network because of privacy and security concerns.
//
// With this subcommand, you can inspect your logs.
// Get is a subcommand of suspicion command where you can query
// suspicions on your agent, without leaking any sensitive data outside of your network.
func GetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Command to query suspicions on your agent",
		Long: `StrixEye agent logs suspicions(anomalies) of your visitors for a span of time.
If you are our customer, then you know that StrixEye agent runs on your internal network and doesn't leak
any data outside of your network/its network because of privacy and security concerns.

With this subcommand, you can inspect your logs.
Get is a subcommand of suspicion command where you can query
suspicions on your agent, without leaking any sensitive data outside of your network.`,
		RunE: getSuspicionCmd,
	}
	
	// declaring local flags used by get suspicion commands.
	getCmd.Flags().StringSliceP(
		"suspects", "s", nil, "Comma separated values of suspect uuids. "+
			"--suspects suspect-id,suspect-id2",
	)
	getCmd.Flags().StringSliceP(
		"suspicions", "a", nil, "Comma separated values of suspicion uuids. "+
			"--suspicions suspicion-id,suspicion-id2",
	)
	getCmd.Flags().StringSliceP(
		"trips", "t", nil, "Comma separated values of trip uuids. "+
			"--trips trip-id,trip-id2",
	)
	
	getCmd.Flags().IntP("limit", "l", 5, "Max number of suspicions you want to be displayed --limit 5")
	
	getCmd.Flags().IntP(
		"since", "i", 0,
		"Queries only suspicions after given time --since [epoch in seconds]  You can get current" +
			" timestamp with date +%s",
	)
	
	return getCmd
}

// getSuspicionCmd implements get suspicion logic.
func getSuspicionCmd(cmd *cobra.Command, _ []string) error {
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
	
	// parse and set list of suspicions to be queried
	suspicionIds, err := cmd.Flags().GetStringSlice("suspicions")
	if err != nil {
		return err
	} else if len(suspicionIds) > 0 {
		queryArgs.SuspicionIds = suspicionIds
	}
	
	// parse and set list of suspicions to be queried
	tripIds, err := cmd.Flags().GetStringSlice("trips")
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
	
	// parse oldest timestamp to be queries
	sinceTime, err := cmd.Flags().GetInt("since")
	if err != nil {
		return err
		
	} else if sinceTime > 0 {
		queryArgs.SinceTime = int64(sinceTime)
	}
	
	// get suspicions with parsed query arguments for this subcommand
	suspicions, err := Get(cliConfig, queryArgs)
	if err != nil {
		return err
		
	}
	
	// marshal result with colors
	data, err := prettyjson.Marshal(suspicions)
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
func Get(cliConfig config.Cli, args QueryArgs) ([]Suspicion, error) {
	var (
		dbConfig config.Database
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

// Get retrieves all suspicions that matches given query args. Check out suspicions.
// QueryArgs for more information about existing filters.
func get(dbConfig config.Database, args QueryArgs) ([]Suspicion, error) {
	var (
		err    error
		db     *gorm.DB
		result []Suspicion
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
	
	// filter by suspicion ids
	if args.SuspicionIds != nil {
		tx = tx.Where(args.SuspectIds)
	}
	
	// filter by suspect ids
	if args.SuspectIds != nil {
		tx = tx.Where("profile_id IN ?", args.SuspectIds)
	}
	
	// filter by suspect ids
	if args.TripsIds != nil {
		tx = tx.Where("trip_id IN ", args.TripsIds)
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
