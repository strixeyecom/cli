package repository

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
 */

// global constants for file
const (
	semVerRegExp = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
)

// global variables (not cool) for this package
var (
	validate = validator.New()
)

// Adding custom validator operators for our usecase
func init() {
	// register custom validation: rfe(Required if Field is Equal to some value).
	err := validate.RegisterValidation(
		`port`, func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			my, err := strconv.Atoi(value)
			if err != nil {
				return false
			}
			if my <= 0 || my >= 65336 {
				return false
			}
			return true
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	// register custom validation: semantic version
	err = validate.RegisterValidation(
		`semver`, func(fl validator.FieldLevel) bool {
			version := fl.Field().String()
			rex, err := regexp.Compile(semVerRegExp)
			if err != nil {
				logrus.Fatal(err)
			}

			// temporary edge case handling
			if version == "staging" || version == "latest" {
				return true
			}

			// return true if field is a semantic version
			version = strings.TrimPrefix(version, "v")
			pass := rex.MatchString(version)

			return pass
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}

// Database stores credentials and configurations about strixeye agent database.
type Database struct {
	DBAddr               string `mapstructure:"DB_ADDR" json:"db_addr"  yaml:"db_addr" flag:"db-addr"`
	DBUser               string `mapstructure:"DB_USER" json:"db_user" validate:"omitempty" yaml:"db_user" flag:"db-user"`
	DBPass               string `mapstructure:"DB_PASS" json:"db_pass" validate:"omitempty" yaml:"db_pass" flag:"db-pass"`
	DBName               string `mapstructure:"DB_NAME" json:"db_name" validate:"omitempty" yaml:"db_name" flag:"db-name"`
	DBPort               string `mapstructure:"DB_PORT" json:"db_port" validate:"port" yaml:"db_port" flag:"db-port"`
	OverrideRemoteConfig bool   `mapstructure:"DB_OVERRIDE" json:"override_remote_config" yaml:"override_remote_config" flag:"override-remote-config"`
	testContainerName    string `flag:"test_container_name"`
}

func (d *Database) TestContainerName() string {
	return d.testContainerName
}

func (d *Database) SetTestContainerName(testContainerName string) {
	d.testContainerName = testContainerName
}

// DSN creates a dsn url from database config. DSN is used to connect to servers,
// this function creates one specific for gorm.
//
// See https://gorm.io/docs/connecting_to_the_database.html
func (d Database) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s", d.DBUser,
		d.DBPass, d.DBAddr, d.DBPort, d.DBName,
	)
}

// Validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (d Database) Validate() error {
	return validate.Struct(d)
}
