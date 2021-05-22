package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

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
const (
	APIUrl       = "https://dashboard.***REMOVED***"
	APITokenName = "Authorization"
)

// global variables (not cool) for this file
var ()

// agentsResponse what user api returns in case of error
type agentsResponse struct {
	Data   []config.AgentInformation `json:"data"`
	Meta   config.Meta               `json:"meta"`
	Status string                    `json:"status"`
}

// GetAgents returns list of agents from user api, parses and validates information.
func GetAgents(apiToken string) ([]config.AgentInformation, error) {
	var (
		err error

		url string

		req  *http.Request
		resp *http.Response
	)

	// create url
	url = fmt.Sprintf("%s/api/agents", APIUrl)

	// create request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// authentication is made via bearer token over http headers.
	tokenHeader := fmt.Sprintf("Bearer %s", apiToken)
	req.Header.Add(APITokenName, tokenHeader)
	req.Header.Add("accept", "application/json")

	// create client to do the request
	client := http.Client{
		Timeout: time.Second * 5,
	}

	// fetch information
	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to user api")
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
