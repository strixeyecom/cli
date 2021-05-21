package suspect

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
	get suspects from your agent database
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Get retrieves all suspects that matches given query args. Check out suspects.
// QueryArgs for more information about existing filters.
func Get(dbConfig config.Database, args QueryArgs) ([]Suspect, error) {
	var (
		err    error
		db     *gorm.DB
		result []Suspect
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

	// filter by score
	if args.Score != 0 {
		tx = tx.Where("score > ? ", args.Score)
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
