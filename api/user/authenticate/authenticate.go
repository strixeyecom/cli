package authenticate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/pkg/errors"
	
	"github.com/strixeyecom/cli/api/user/repository"
	`github.com/strixeyecom/cli/domain/cli`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	Authenticate users over StrixEye User API.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// Authenticate checks if given user api token and url is valid.
// it returns error for status codes other than 200.
func Authenticate(cliConfig cli.Cli) error {
	return authenticate(cliConfig.UserAPIToken, cliConfig.APIDomain)
}

// authenticate checks if given user api key and url is valid.
// it returns error for status codes other than 200.
func authenticate(apiToken, apiDomain string) error {
	var (
		err  error
		resp *http.Response
	)

	resp, err = repository.UserAPIRequest(http.MethodGet, "", nil, apiToken, apiDomain)
	if err != nil {
		return errors.Wrap(err, "failed to complete user api request to agents")
	}
	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "bad response body")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// handle reject/fail responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"sorry, please double check your credentials. "+
				"Status Code : %d, error message : %s", resp.StatusCode, body,
		)
	}

	return nil
}
