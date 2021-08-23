package suspicion

import (
	`github.com/fatih/color`
	`github.com/hokaccha/go-prettyjson`
	"github.com/pkg/errors"
	`github.com/spf13/cobra`
	`github.com/spf13/viper`
	userconfig `github.com/strixeyecom/cli/api/user/agent`
	models `github.com/strixeyecom/cli/domain/repository`
	"gorm.io/gorm"
	
	`github.com/strixeyecom/cli/cli/command/repository`
	`github.com/strixeyecom/cli/domain/cli`
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
		"Queries only suspicions after given time --since [epoch in seconds]  You can get current"+
			" timestamp with date +%s",
	)
	
	return getCmd
}

// getSuspicionCmd implements get suspicion logic.
func getSuspicionCmd(cmd *cobra.Command, _ []string) error {
	var (
		cliConfig cli.Cli
		err       error
	)
	
	// get cli config for authentication
	err = viper.Unmarshal(&cliConfig)
	if err != nil {
		return err
	}
	
	queryArgs := models.SuspicionQueryArgs{Limit: 1}
	
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
	
	// Insert extra information from StrixEye API
	for i, suspicion := range suspicions {
		suspicions[i].Domain, err = GetDomainInformation(cliConfig, suspicion.DomainId)
		if err != nil {
			return errors.Wrap(err, "can not fetch domain information")
		}
		
		tmp, err := trip.GetTrips(
			cliConfig,
			models.TripQueryArgs{TripsIds: []string{suspicion.TripId}, Verbose: true, Limit: 1},
		)
		if err != nil {
			return errors.Wrap(err, "can not extract trip information")
		}
		
		suspicions[i].Trip = tmp[0]
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
func Get(cliConfig cli.Cli, args models.SuspicionQueryArgs) ([]models.Suspicion, error) {
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

// Get retrieves all suspicions that matches given query args. Check out suspicions.
// QueryArgs for more information about existing filters.
func get(dbConfig models.Database, args models.SuspicionQueryArgs) ([]models.Suspicion, error) {
	var (
		err    error
		db     *gorm.DB
		result []models.Suspicion
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

// GetDomainInformation returns domain information
func GetDomainInformation(cliConfig cli.Cli, domainID string) (models.Domain, error) {
	return getDomain(cliConfig.UserAPIToken, cliConfig.APIDomain, domainID)
}

// getDomain returns list of agents from user api, parses and validates information.
func getDomain(apiToken, apiDomain, domainID string) (models.Domain, error) {
	var (
		err  error
		resp *http.Response
	)
	if domainID == "" {
		return models.Domain{}, errors.New("no domain id given")
	}
	url := fmt.Sprintf("/domains/%s", domainID)
	resp, err = repository2.UserAPIRequest(http.MethodGet, url, nil, apiToken, apiDomain)
	
	if err != nil {
		return models.Domain{}, errors.Wrap(err, "failed to complete user api request to agents")
	}
	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Domain{}, errors.Wrap(err, "bad response body")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	
	// handle reject/fail responses
	if resp.StatusCode != http.StatusOK {
		return models.Domain{}, fmt.Errorf(
			"sorry, please double check your credentials. "+
				"Status Code : %d, error message : %s", resp.StatusCode, body,
		)
	}
	
	// if status is ok, than this is possibly a api success response
	var apiResponse models.DomainMessage
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return models.Domain{}, errors.Wrap(
			err,
			"api says response is okay but possibly there is a misunderstanding",
		)
	}
	
	return apiResponse.Data, nil
}
