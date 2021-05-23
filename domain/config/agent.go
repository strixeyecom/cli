package config

import "fmt"

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	agent.go stores structural data for agent's config
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Agent config is a very simple but important config.
//
// It is used by strixeyed to keep agent token,
// so that it can authenticate itself and establish a connection to connector.
//
// We don't see this struct being used too often.
type Agent struct {
	Versions  Versions  `json:"versions"`
	Auth      auth      `json:"auth"`
	Addresses addresses `json:"addresses"`
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

// AgentInformation keeps all information relevant to an agent instance.
type AgentInformation struct {
	ID        string      `json:"id"`
	CompanyID string      `json:"company_id"`
	Name      string      `json:"name"`
	IPAddress string      `json:"ip_address"`
	Config    stackConfig `json:"config"`
	Domains   []domains   `json:"domains"`
}

func (a AgentInformation) String() string {
	return fmt.Sprintf(
		"Name: %s,  Id : %s , IP : %s", a.Name, a.ID,
		a.IPAddress,
	)
}
