package config

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	versions.go stores structural data for version messages.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Versions stores latest version informations which then used to:
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

func (m Versions) Validate() error {
	err :=  validate.Struct(m)
	return err
}
