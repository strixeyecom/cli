package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	
	"github.com/pkg/errors"
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	stack.go handles how we interact with stack config.
*/

// global constants for file
const ()

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

// stackConfig is THE most important config in all strixeye universe. This is how agents self-update,
// self-deploy themselves, switch from kubernetes to docker, change database user,
// send data more frequently, send data less frequently, use nginx, use apache, use http,
// use https and all those stuff. fun.
type stackConfig struct {
	Addresses  addresses `json:"addresses"`
	UseHTTPS   bool      `json:"use_https"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Deployment string    `json:"deployment"`
	Database   Database  `json:"database"`
	Broker     broker    `json:"broker"`
	Scheduler  scheduler `json:"scheduler"`
	Engine     engine    `json:"engine"`
	Sensor     sensor    `json:"sensor"`
	Profiler   profiler  `json:"profiler"`
	Intervals  intervals `json:"intervals"`
	Paths      paths     `json:"paths"`
}

// Save stores stackConfig as json to a given path.
func (config stackConfig) Save(filePath string) error {
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
func (config stackConfig) Validate() error {
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

// addresses keep where/how agent connects to
type addresses struct {
	// Scheme means here that whether it is a websocket connection, and therefore "ws" or "wss"
	// or a normal http connection, "http" or "https"
	ConnectorScheme string `json:"connector_scheme"`
	
	// ConnectorAddress keeps connector's location.
	// This is usually fixed since strixeye has a cloud management panel and it is on a predefined domain.
	// But still, it has a hostname validation
	ConnectorAddress string `json:"connector_address" validate:"hostname"`
	
	// ConnectorPort has same explanation with ConnectorAddress field of the same struct.
	ConnectorPort string `json:"connector_port" validate:"port"`
}

// validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (a addresses) validate() error {
	if a.ConnectorScheme != "wss" && a.ConnectorScheme != "ws" {
		return errors.New("bad connector scheme")
	}
	
	return validate.Struct(a)
}

// Database stores credentials and configurations about strixeye agent database.
type Database struct {
	DBAddr               string `mapstructure:"DB_ADDR" json:"db_addr" validate:"hostname"`
	DBUser               string `mapstructure:"DB_USER" json:"db_user" validate:"omitempty"`
	DBPass               string `mapstructure:"DB_PASS" json:"db_pass" validate:"omitempty"`
	DBName               string `mapstructure:"DB_NAME" json:"db_name" validate:"omitempty"`
	DBPort               string `mapstructure:"DB_PORT" json:"db_port" validate:"port"`
	OverrideRemoteConfig bool   `mapstructure:"DB_OVERRIDE" json:"override_remote_config"`
}

// DSN creates a dsn url from database config. DSN is used to connect to servers,
// this function creates one specific for gorm.
//
// See https://gorm.io/docs/connecting_to_the_database.html
func (d Database) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", d.DBUser,
		d.DBPass, d.DBAddr, d.DBPort, d.DBName,
	)
}

// Validate checks for the fields of given instance.
// check for struct type definition for more documentation about fields and their validation functions.
func (d Database) Validate() error {
	return validate.Struct(d)
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

type domains struct {
	ID        string      `json:"id"`
	CompanyID string      `json:"company_id"`
	Domain    string      `json:"domain"`
	DeletedAt interface{} `json:"deleted_at"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Pivot     pivot       `json:"pivot"`
}
