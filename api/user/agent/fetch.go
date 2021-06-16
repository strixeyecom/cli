package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	repository2 `github.com/usestrix/cli/domain/repository`
	
	"github.com/usestrix/cli/api/user/repository"
	"github.com/usestrix/cli/domain/agent"
	"github.com/usestrix/cli/domain/cli"
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	fetch configuration of current agent from strixeye user api
*/

// global constants for file
const (
	APITokenName = "Authorization"
)

// global variables (not cool) for this file
var ()

// agentResponse what user api returns in case of error
type agentResponse struct {
	Data   agent.AgentInformation `json:"data"`
	Status string                 `json:"status"`
}

// GetAgentConfig return stack configuration for given agent.
func GetAgentConfig(cliConfig cli.Cli) (agent.AgentInformation, error) {
	return getAgent(cliConfig.UserAPIToken, cliConfig.APIUrl, cliConfig.CurrentAgentID)
}

// getAgents returns list of agents from user api, parses and validates information.
func getAgent(apiToken, apiURL, agentID string) (agent.AgentInformation, error) {
	var (
		err  error
		resp *http.Response
	)

	url := fmt.Sprintf("/agents/%s", agentID)
	resp, err = repository.UserAPIRequest(http.MethodGet, url, nil, apiToken, apiURL)

	if err != nil {
		return agent.AgentInformation{}, errors.Wrap(err, "failed to complete user api request to agents")
	}
	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return agent.AgentInformation{}, errors.Wrap(err, "bad response body")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// handle reject/fail responses
	if resp.StatusCode != http.StatusOK {
		return agent.AgentInformation{}, fmt.Errorf(
			"sorry, please double check your credentials. "+
				"Status Code : %d, error message : %s", resp.StatusCode, body,
		)
	}

	// if status is ok, than this is possibly a api success response
	var apiResponse agentResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return agent.AgentInformation{}, errors.Wrap(
			err,
			"api says response is okay but possibly there is a misunderstanding",
		)
	}

	return apiResponse.Data, nil
}

// GetLatestVersions return latest image/binary versions
func GetLatestVersions(cliConfig cli.Cli) (repository2.Versions, error) {
	return getVersions(cliConfig.APIUrl)
}

// getVersions returns latest image versions you can use with StrixEye Agents.
func getVersions(apiURL string) (repository2.Versions, error) {
	var (
		err  error
		resp *http.Response
	)

	url := fmt.Sprintf("%s/versions", apiURL)
	resp, err = http.Get(url)

	if err != nil {
		return repository2.Versions{}, errors.Wrap(err, "failed to complete request to get versions")
	}
	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return repository2.Versions{}, errors.Wrap(err, "bad response body")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// handle reject/fail responses
	if resp.StatusCode != http.StatusOK {
		return repository2.Versions{}, fmt.Errorf(
			"sorry, please double check your credentials. "+
				"Status Code : %d, error message : %s", resp.StatusCode, body,
		)
	}

	// if status is ok, than this is possibly a api success response
	var apiResponse repository2.APIVersionsMessage
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return repository2.Versions{}, errors.Wrap(
			err,
			"api says response is okay but possibly there is a misunderstanding",
		)
	}

	if len(apiResponse.Data) == 0 {
		return repository2.Versions{}, errors.Wrap(
			err,
			"under maintenance, please try again in a few minutes",
		)
	}

	// turn api response to usable version structure
	versions, err := apiResponse.ToVersions()
	if err != nil {
		return repository2.Versions{}, err
	}
	return versions, nil
}
