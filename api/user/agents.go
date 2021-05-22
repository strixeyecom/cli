package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/usestrix/cli/api/user/repository"
	"github.com/usestrix/cli/domain/config"
)

/*
	Created by aomerk at 5/22/21 for project cli
*/

/*
	get basic information about agent
*/

// global constants for file

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// agentsResponse what user api returns in case of error
type agentsResponse struct {
	Data   []config.AgentInformation `json:"data"`
	Meta   config.Meta               `json:"meta"`
	Status string                    `json:"status"`
}

// GetAgents returns list of agents from user api, parses and validates information.
func GetAgents(cliConfig config.Cli) ([]config.AgentInformation, error) {
	return getAgents(cliConfig.UserAPIToken, cliConfig.APIUrl)
}

// getAgents returns list of agents from user api, parses and validates information.
func getAgents(apiToken, apiURL string) ([]config.AgentInformation, error) {
	var (
		err  error
		resp *http.Response
	)

	resp, err = repository.UserAPIRequest(http.MethodGet, "/agents", nil, apiToken, apiURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to complete user api request to agents")
	}
	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "bad response body")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// handle reject/fail responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"sorry, please double check your credentials. "+
				"Status Code : %d, error message : %s", resp.StatusCode, body,
		)
	}

	// if status is ok, than this is possibly a api success response
	var apiResponse agentsResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"api says response is okay but possibly there is a misunderstanding",
		)
	}

	return apiResponse.Data, nil
}
