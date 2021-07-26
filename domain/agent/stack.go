package agent

import (
	"encoding/json"
	"io/ioutil"
	`regexp`
	`strconv`
	`strings`
	"time"
	
	`github.com/go-playground/validator`
	"github.com/pkg/errors"
	`github.com/sirupsen/logrus`
	
	`github.com/usestrix/cli/domain/repository`
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	stack.go handles how we interact with stack config.
*/
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


// global constants for file
const (
	semVerRegExp = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
)
// APIStackResponse is what we usually get from user api as response when we try to retrieve an agent's
// stack config.
// there are two main cases. We get stack information,
// or we get an error. If we get stack information, it means that the status is "ok", if not "ok",
// it is some error that we can simply print out. ~~~~simply print out~~~~~ We
//
// We can't simply print out. Two types of responses are just not compatible!
//
// 200 and >300 response codes must be handled differently.
type APIStackResponse struct {
	// Status is usually "ok" or "failed"
	Status string `json:"status"`
	
	// Data is usually keeps the error message or out stack config
	Stack AgentInformation `json:"data"`
}

// APIErrorResponse is usually how user api returns error responses. Of course,
// saying it is a map[string]interface{} is not "how it returns"
//
// But it is as close as we can get with static typing.
type APIErrorResponse struct {
	// Status is usually "ok" or "failed"
	Status string `json:"status"`
	
	// Data is usually keeps the error message or out stack config
	Stack map[string]interface{} `json:"data"`
}

// StackConfig is THE most important config in all strixeye universe. This is how agents self-update,
// self-deploy themselves, switch from kubernetes to docker, change database user,
// send data more frequently, send data less frequently, use nginx, use apache, use http,
// use https and all those stuff. fun.
type StackConfig struct {
	Addresses  Addresses           `json:"addresses"`
	UseHTTPS   bool                `json:"use_https"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	Deployment string              `json:"deployment"`
	Database   repository.Database `json:"database"`
	Broker     broker              `json:"broker"`
	Scheduler  scheduler           `json:"scheduler"`
	Engine     engine              `json:"engine"`
	Sensor     sensor              `json:"sensor"`
	Profiler   profiler            `json:"profiler"`
	Intervals  intervals `json:"intervals"`
	Paths      paths     `json:"paths"`
}

// Save stores stackConfig as json to a given path.
func (config StackConfig) Save(filePath string) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	
	// using default file permission and writing to path
	err = ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		return err
	}
	
	return nil
}

// Validate validates incoming stack config, returns error if unexpected
func (config StackConfig) Validate() error {
	var err error
	err = config.Broker.validate()
	if err != nil {
		return err
	}
	
	err = config.Intervals.validate()
	if err != nil {
		return err
	}
	
	err = config.Scheduler.validate()
	if err != nil {
		return err
	}
	
	err = config.Profiler.validate()
	if err != nil {
		return err
	}
	
	err = config.Engine.validate()
	if err != nil {
		return err
	}
	
	err = config.Sensor.validate()
	if err != nil {
		return err
	}
	
	err = config.Addresses.validate()
	if err != nil {
		return err
	}
	
	err = config.Database.Validate()
	if err != nil {
		return err
	}
	
	return err
}

// Addresses keep where/how agent connects to
type Addresses struct {
	// Scheme means here that whether it is a websocket connection, and therefore "ws" or "wss"
	// or a normal http connection, "http" or "https"
	ConnectorScheme string `json:"connector_scheme"`
	
	// ConnectorAddress keeps connector's location.
	// This is usually fixed since strixeye has a cloud management panel and it is on a predefined domain.
	// But still, it has a hostname validation
	ConnectorAddress string `json:"connector_address"`
	
	// ConnectorPort has same explanation with ConnectorAddress field of the same struct.
	ConnectorPort string `json:"connector_port" validate:"port"`
}

// validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (a Addresses) validate() error {
	if a.ConnectorScheme != "wss" && a.ConnectorScheme != "ws" {
		return errors.New("bad connector scheme")
	}
	
	return validate.Struct(a)
}

type broker struct {
	BrokerHostname string `json:"broker_hostname" validate:"hostname"`
	BrokerUsername string `json:"broker_username" validate:"omitempty"`
	BrokerPrefix   string `json:"broker_prefix"`
	BrokerPassword string `json:"broker_password" validate:"omitempty"`
	BrokerListen   string `json:"broker_port" validate:"port"`
}

// validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (b broker) validate() error {
	if b.BrokerPrefix != "amqp" {
		return errors.New("bad broker uri prefix")
	}
	return validate.Struct(b)
}

type scheduler struct {
	SchedulerListen string `json:"scheduler_listen" validate:"port"`
}

func (s scheduler) validate() error {
	return validate.Struct(s)
}

type engine struct {
	Address      string `json:"address" validate:"hostname"`
	EngineListen string `json:"engine_listen" validate:"port"`
}

// validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (e engine) validate() error {
	return validate.Struct(e)
}

type sensor struct {
	IntegrationName string `json:"integration_name"`
	SensorListen    string `json:"sensor_listen" validate:"port"`
}

// validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (s sensor) validate() error {
	if s.IntegrationName != "nginx" && s.IntegrationName != "apache" {
		return errors.New("bad integration name in configuration")
	}
	return validate.Struct(s)
}

type profiler struct {
	ProfilerListen string `json:"profiler_listen" validate:"port"`
}

func (p profiler) validate() error {
	return validate.Struct(p)
}

type intervals struct {
	SystemStatsIntervalSecond int `json:"system_stats_interval_second" validate:"gt=0"`
}

// validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (i intervals) validate() error {
	return validate.Struct(i)
}

type tlsKeys struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

type paths struct {
	KubeConfig string  `json:"kube_config"`
	TLSKeys    tlsKeys `json:"tls_keys"`
}

type pivot struct {
	AgentID  string `json:"agent_id"`
	DomainID string `json:"domain_id"`
}

type Domains struct {
	ID        string      `json:"id"`
	CompanyID string      `json:"company_id"`
	Domain    string      `json:"domain"`
	DeletedAt interface{} `json:"deleted_at"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Pivot     pivot       `json:"pivot"`
}
