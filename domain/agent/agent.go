package agent

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/usestrix/cli/domain/repository"
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	agent.go stores structural data for agent's config
*/

// global constants for file
const (
	StatusOK = "ok"
)

// global variables (not cool) for this file
var ()

// Agent config is a very simple but important config.
//
// It is used by strixeyed to keep agent token,
// so that it can authenticate itself and establish a connection to connector.
//
// We don't see this struct being used too often.
type Agent struct {
	Versions  repository.Versions `json:"versions"`
	Auth      auth                `json:"auth"`
	Addresses addresses           `json:"addresses"`
}

type auth struct {
	CompanyID    string `json:"company_id"`
	CompanyToken string `json:"company_token"`
	AgentID      string `json:"agent_id"`
}

// Meta keeps meta information about a companies all agents. Not very important.
type Meta struct {
	Total       int         `json:"total"`
	Count       int         `json:"count"`
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	NextPageURL interface{} `json:"next_page_url"`
	PrevPageURL interface{} `json:"prev_page_url"`
}

// Versions keeps all version information for StrixEye stack.
type Versions struct {
	Manager   string `json:"manager_version"`
	Database  string `json:"database_version"`
	Engine    string `json:"engine_version"`
	Profiler  string `json:"profiler_version"`
	Queue     string `json:"queue_version"`
	Scheduler string `json:"scheduler_version"`
	Sensor    string `json:"sensor_version"`
	Installer string `json:"installer_version"`
}

// AgentInformation keeps all information relevant to an agent instance.
type AgentInformation struct {
	ID        string      `json:"id"`
	CompanyID string      `json:"company_id"`
	Name      string      `json:"name"`
	IPAddress string      `json:"ip_address"`
	Config    StackConfig `json:"config"`
	Domains   []Domains   `json:"domains"`
}

func (a AgentInformation) String() string {
	return fmt.Sprintf(
		"Name: %s,\tId: %s\tIP: %s", a.Name, a.ID,
		a.IPAddress,
	)
}

// APIVersionsMessage is how versions endpoint returns. Here,
// instead of implementing Marshaller and Unmarshaler interfaces,
// I am simply using a whole different struct because the way the API returns is a bit silly in terms of
// static typed languages like golang.
type APIVersionsMessage struct {
	Status string `json:"status"`
	Data   []struct {
		Key       string      `json:"key"`
		Value     string      `json:"value"`
		CreatedAt interface{} `json:"created_at"`
		UpdatedAt *time.Time  `json:"updated_at"`
		DeletedAt interface{} `json:"deleted_at"`
	} `json:"data"`
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
func decode(s APIVersionsMessage, i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("decode requires non-nil pointer")
	}
	// get the value that the pointer v points to.
	v = v.Elem()
	// assume that the input is valid.
	for _, kv := range s.Data {
		// strip version suffix
		fieldName := strings.ReplaceAll(kv.Key, "_version", "")

		// make first letter capital like an exported field name
		fieldName = strings.Title(fieldName)

		// get field
		f := v.FieldByName(fieldName)

		// make sure that this field is defined, and can be changed.
		if !f.IsValid() || !f.CanSet() {
			continue
		}

		if kv.Value == "" {
			return errors.New("sorry, api returned a bad versions message")
		}
		// assume all the fields are type string.
		f.SetString(kv.Value)
	}
	return nil
}
