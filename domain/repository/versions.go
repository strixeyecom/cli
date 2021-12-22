package repository

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	versions.go stores structural data for version messages.
*/

// global constants for file.
const (
	// StatusOK is status message from StrixEye API for 2xx messages.
	StatusOK = "ok"
)

// global variables (not cool) for this file.
var ()

// Versions stores latest version information which then used to:
//
// - pull stack container images
// - download installer
// - download strixeyed
//
// Versions keeps all version information for StrixEye stack.
type Versions struct {
	Manager   Version `json:"manager_version" image:"-"`
	Database  Version `json:"database_version" image:"database"`
	Engine    Version `json:"engine_version" image:"engine"`
	Profiler  Version `json:"profiler_version" image:"profiler"`
	Queue     Version `json:"queue_version" image:"queue"`
	Scheduler Version `json:"scheduler_version" image:"scheduler"`
	Sensor    Version `json:"sensor_version" image:"sensor"`
	Installer Version `json:"installer_version" image:"-"`
}

func (v *Versions) VisitAll(g func(value reflect.StructField) error) error {
	dst := reflect.ValueOf(v)
	if dst.Kind() != reflect.Ptr || dst.IsNil() {
		return fmt.Errorf("decode requires non-nil pointer")
	}

	// get the value that the pointer dst points to.
	dst = dst.Elem()
	numOfFields := dst.NumField()
	for i := 0; i < numOfFields; i++ {
		field := dst.Type().Field(i)

		err := g(field)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate returns error if there is a problem with validation.
// Check struct type definition of Versions for more information
func (v *Versions) Validate() error {
	err := validate.Struct(v)
	return err
}

func (message APIVersionsMessage) ToVersions() (Versions, error) {
	var (
		versionsResult = Versions{}
	)

	if message.Status != StatusOK {
		return versionsResult, errors.New("API response is not successful")
	}

	// iterate over fields and fill your result
	// there are also libraries for JSON where you can directly access data from json,
	// but this way, it is easier to detect mistakes, yet a bit worse to maintain.
	err := decode(message, &versionsResult)
	if err != nil {
		return Versions{}, err
	}

	return versionsResult, nil
}

// decode fills versions message using reflect package.
// Simple switch factory would be more than enough but since this is a really boring thing to do,
// I wanted to spice things up a little.
func decode(s APIVersionsMessage, i interface{}) error {
	dst := reflect.ValueOf(i)
	if dst.Kind() != reflect.Ptr || dst.IsNil() {
		return fmt.Errorf("decode requires non-nil pointer")
	}
	// get the value that the pointer dst points to.
	dst = dst.Elem()
	// assume that the input is valid.
	for _, kv := range s.Data {
		// strip version suffix
		fieldName := strings.ReplaceAll(kv.Key, "_version", "")

		// make first letter capital like an exported field name
		fieldName = strings.Title(fieldName)

		// get field
		f := dst.FieldByName(fieldName)

		// make sure that this field is defined, and can be changed.
		if !f.IsValid() || !f.CanSet() {
			continue
		}

		if kv.Value.Version == "" {
			return errors.New("sorry, api returned a bad versions message")
		}

		f.Set(reflect.ValueOf(kv.Value))
	}
	return nil
}

// APIVersionsMessage is how versions endpoint returns. Here,
// instead of implementing Marshaller and Unmarshaler interfaces,
// I am simply using a whole different struct because the way the API returns is a bit silly in terms of
// static typed languages like golang.
type APIVersionsMessage struct {
	Status string `json:"status"`
	Data   []struct {
		Key       string  `json:"key"`
		Value     Version `json:"value"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
		DeletedAt string  `json:"deleted_at"`
	} `json:"data"`
}

// Version keeps version/download related information for images and binaries.
type Version struct {
	Checksum string `json:"checksum"`
	Size     int64  `json:"size"`
	Version  string `json:"version" validate:"semver"`
}
