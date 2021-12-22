package suspect

import (
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	userconfig "github.com/strixeyecom/cli/api/user/agent"
	"github.com/strixeyecom/cli/cli/command/repository"
	"github.com/strixeyecom/cli/domain/cli"
	models "github.com/strixeyecom/cli/domain/repository"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	get suspects from your agent database
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// GetCommand returns cli command to query suspects.
// Command to query suspects on your agent
// StrixEye agent logs suspicions of your visitors for a span of time.
// If you are our customer, then you know that StrixEye agent runs on your internal network and doesn't leak
// any data outside of your network/its network because of privacy and security concerns.
//
// With this subcommand, you can inspect your logs get is a subcommand of trip command where you can query
// trips on your agent, without leaking any sensitive data outside of your network.
func GetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Command to query suspects on your agent",
		Long: `Command to query suspects on your agent
StrixEye agent logs requests/responses of your visitors for a span of time.
If you are our customer, then you know that StrixEye agent runs on your internal network and doesn't leak
any data outside of your network/its network because of privacy and security concerns.

With this subcommand, you can inspect your logs get is a subcommand of trip command where you can query
trips on your agent, without leaking any sensitive data outside of your network`,
		RunE: getSuspectCmd,
	}

	// declaring local flags used by get trip commands.
	getCmd.Flags().StringSliceP(
		"ids", "i", nil, "Comma separated values of suspect uuids. "+
			"--ids suspect-id,suspect-id2",
	)

	getCmd.Flags().IntP("limit", "l", 5, "Max number of suspects you want to be displayed --limit 5")

	getCmd.Flags().BoolP(
		"verbose", "v", true, "To hide field values like headers, say --verbose=false",
	)
	getCmd.Flags().IntP(
		"since", "s", 0,
		"Queries only suspects detected after given time --since [epoch in seconds]  You can get current"+
			" timestamp with date +%s",
	)

	getCmd.Flags().Float64P(
		"min-score", "m", 0,
		"Queries only suspects with scores bigger than X --min-score [0-100] Default value is 0.  ",
	)

	return getCmd
}

// getSuspectCmd implements GetCommand logic.
func getSuspectCmd(cmd *cobra.Command, _ []string) error {
	var (
		cliConfig cli.Cli
		err       error
	)

	// get cli config for authentication
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}

	queryArgs := models.SuspectQueryArgs{Limit: 1}

	// parse and set list of suspects to be queried
	suspects, err := cmd.Flags().GetStringSlice("ids")
	if err != nil {
		return err
	} else if len(suspects) > 0 {
		queryArgs.SuspectIds = suspects
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

	// parse minimum score of suspects
	minScore, err := cmd.Flags().GetFloat64("min-score")
	if err != nil {
		return err
	} else if minScore >= 0 {
		queryArgs.MinScore = minScore
	}

	// get trips with parsed query arguments for this subcommand
	trips, err := Get(cliConfig, queryArgs)
	if err != nil {
		return err

	}

	// print out query settings
	color.Blue(queryArgs.String())

	// print out result
	_, err = pp.Print(trips)
	if err != nil {
		return err
	}

	return nil
}

// Get is a temporary method to satisfy the authentication process.
func Get(cliConfig cli.Cli, args models.SuspectQueryArgs) ([]models.Suspect, error) {
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

// get retrieves all suspects that matches given query args. Check out suspects.
// QueryArgs for more information about existing filters.
func get(dbConfig models.Database, args models.SuspectQueryArgs) ([]models.Suspect, error) {
	var (
		err    error
		db     *gorm.DB
		result []models.Suspect
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

	// preload only first level for now
	tx = tx.Preload(clause.Associations)
	tx = tx.Preload("Trips.StaticChecks")
	tx = tx.Preload("Trips.Request")
	tx = tx.Preload("Trips.Request.Header")

	// filter by score
	if args.MinScore != 0 {
		tx = tx.Where("score > ? ", args.MinScore)
	}

	// filter by suspect ids
	if args.SuspectIds != nil {
		tx = tx.Where(args.SuspectIds)
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
