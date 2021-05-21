package config

import (
	"encoding/json"

	"github.com/pkg/errors"
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	versions.go stores structural data for version messages.
*/

// global constants for file.
const ()

// global variables (not cool) for this file.
var ()

// FromAPIResponse is a tricky pattern. Using a single Unmarshal method for both Versions itself and
// versionsJSON sounded good but had a bad implementation.
func (v *Versions) FromAPIResponse(apiVersions versionsJSON) error {
	var (
		err error
	)

	// Iterate over fields and try to find something we need as version
	for _, datum := range apiVersions.Data {
		switch datum.Key {
		case "sensor_version":
			v.SensorVersion = datum.Value
		case "queue_version":
			v.QueueVersion = datum.Value
		case "profiler_version":
			v.ProfilerVersion = datum.Value
		case "scheduler_version":
			v.SchedulerVersion = datum.Value
		case "engine_version":
			v.EngineVersion = datum.Value
		case "manager_version":
			v.ManagerVersion = datum.Value
		case "database_version":
			v.DatabaseVersion = datum.Value
		}
	}

	// it is more than a unmarshaler, it is also a validator!!
	err = v.Validate()
	if err != nil {
		return err
	}

	return nil
}

// FromRawApiResponse is a tricky pattern. Using a single Unmarshal method for both Versions itself and
// versionsJSON sounded good but had a bad implementation.
func (v Versions) FromRawAPIResponse(rawData []byte) error {
	var (
		apiVersions versionsJSON
		err         error
	)

	// 	unmarshal to api response. It can have hundreds of different versions,
	// 	even the version of linux kernel version of build machine.
	err = json.Unmarshal(rawData, &apiVersions)
	if err != nil {
		return err
	}

	err = v.FromAPIResponse(apiVersions)
	if err != nil {
		return errors.Wrap(err, "can not parse user api response")
	}
	return nil
}

// Versions stores latest version information which then used to:
//
// - pull stack container images
// - download installer
// - download strixeyed
type Versions struct {
	ManagerVersion   string `json:"manager_version" validate:"semver"`
	DatabaseVersion  string `json:"database_version" validate:"semver"`
	EngineVersion    string `json:"engine_version" validate:"semver"`
	ProfilerVersion  string `json:"profiler_version" validate:"semver"`
	QueueVersion     string `json:"queue_version" validate:"semver"`
	SchedulerVersion string `json:"scheduler_version" validate:"semver"`
	SensorVersion    string `json:"sensor_version" validate:"semver"`
}

// Validate returns error if there is a problem with validation.
// Check struct type definition of Versions for more information
func (v Versions) Validate() error {
	err := validate.Struct(v)
	return err
}

// versionsJSON keeps the version data in a different way,
// because this is how user api models version data.
// It is a bit tricky to extract fields for golang but not really a huge problem.
type versionsJSON struct {
	Status string             `json:"status"`
	Data   []apiSingleVersion `json:"data"`
}

// apiSingleVersion is a generic KV object structure. User api stores a single version in this manner.
type apiSingleVersion struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
