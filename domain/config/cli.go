package config

import (
	`encoding/json`
	`io/ioutil`
	
	`github.com/pkg/errors`
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	Here we keep the necessary configuration of our cli.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Cli keeps needed information to run strixeye cli.
// It's usually kept under the same folder with strixeyed config file if there is one.
//
// It's usually kept in json format.
type Cli struct {
	
	// UserAPIToken is used for the authentication process to Strixeye User API.
	// This api is open to all our customers and feel free to check out the documentation.
	//
	// UserAPIToken is generally sent as Authentication Bearer token over https.
	UserAPIToken string
	
	// strixeye cli is usually designed to be used for a single agent instance at once,
	// while not necessary, this field can be use to save user preference.
	//
	// However, most functions are agent id dependent.
	CurrentAgentId string
}

// Save stores current cli config to given file in json format.
func (c *Cli) Save(filePath string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "can not save cli config to file")
	}
	
	// using rw|rw|rw for permission.
	return ioutil.WriteFile(filePath, data, 0666)
}

// Validate checks for bad/empty input inside config instances. Ids and Tokens are mostly generated by uuids.
func (c *Cli) Validate() error {
	// TODO add validation support
	if c.UserAPIToken == "" || c.CurrentAgentId == "" {
		return errors.New("has empty field")
	}
	return nil
}
