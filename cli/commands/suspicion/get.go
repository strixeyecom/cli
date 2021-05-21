package suspicion

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
	get suspicions from your agent database
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Get retrieves all suspicions that matches given query args. Check out suspicions.
// QueryArgs for more information about existing filters.
func Get(dbConfig config.Database, args QueryArgs) ([]Suspicion, error) {
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

	// preload only first level for now
	tx = tx.Preload(clause.Associations)

	// filter by suspicion ids
	if args.SuspicionIds != nil {
		tx = tx.Where(args.SuspicionIds)
	}

	// filter by suspect ids
	if args.SuspectIds != nil {
		tx = tx.Where("profile_id IN ?", args.SuspectIds)
	}

	// filter by trip ids
	if args.TripsIds != nil {
		tx = tx.Where("profile_id IN ?", args.TripsIds)
	}

	// find all suspicions that matches args criteria
	tx = tx.Find(&result)
	return result, nil
}
