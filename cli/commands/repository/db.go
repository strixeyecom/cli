package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	repository for database controls
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// ConnectToAgentDB establishes a live connection to your agents database.
// Make sure to have permissions and network configurations so that use can connect to database. Usually,
// database ports and hosts are not public in enterprise networks. So, that part is on you to check.
func ConnectToAgentDB(dbConfig config.Database) (*gorm.DB, error) {
	// establish connection
	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN:                       dbConfig.DSN(), // data source name
				DefaultStringSize:         256,            // default size for string fields
				DisableDatetimePrecision:  true,           // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,           // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,           // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false,          // auto configure based on currently MySQL version
			},
		), &gorm.Config{},
	)

	return db, err
}
