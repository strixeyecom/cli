package repository

import (
	"log"
	"os"
	"time"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	`github.com/usestrix/cli/domain/repository`
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
func ConnectToAgentDB(dbConfig repository.Database) (*gorm.DB, error) {
	// orm logger config
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

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
		), &gorm.Config{
			Logger: newLogger,
		},
	)

	return db, err
}
