package cli

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

/*
	Created by aomerk at 5/20/21 for project strixeye
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
			if my <= 1024 || my >= 65336 {
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

// Config is base interface to implement for strixeye configuration structs.
type Config interface {
	// Since most of the config is crucial, validation process is highly encouraged.
	Validate() error

	// Configs are mostly kept as files.
	Save(filePath string) error
}
