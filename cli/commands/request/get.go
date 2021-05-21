package trip

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

// Get retrieves all trips that matches given query args. Check out trips.
// QueryArgs for more information about existing filters.
func Get(dbConfig config.Database, args QueryArgs) ([]Trip, error) {
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
		tx = tx.Where("profile_id = ?", args.SuspectIds)
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
