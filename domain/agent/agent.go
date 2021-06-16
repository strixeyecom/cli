package agent

import (
	"fmt"
	
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
	// StatusOK is status message from StrixEye API for 2xx messages.
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
	Auth      Auth                `json:"auth"`
	Addresses Addresses           `json:"addresses"`
}

// Auth stores authentication credentials of agent.
type Auth struct {
	AgentToken string `json:"agent_token"`
	AgentID    string `json:"agent_id"`
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

// AgentInformation keeps all information relevant to an agent instance.
type AgentInformation struct {
	ID        string      `json:"id"`
	Token     string      `json:"token"`
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
