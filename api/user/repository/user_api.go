package repository

import (
	"fmt"
	"io"
	"net/http"
	"time"
	
	"github.com/pkg/errors"
)

/*
	Created by aomerk at 5/22/21 for project cli
*/

/*
	helper for creating user api requests
*/

// global constants for file
const (
	APITokenName = "Authorization"
)

// global variables (not cool) for this file
var ()

func UserAPIRequest(method, endpoint string, body io.Reader, apiToken, apiDomain string) (
	*http.Response, error,
) {
	var (
		err  error
		url  string
		req  *http.Request
		resp *http.Response
	)
	
	// create url
	url = fmt.Sprintf("https://%s%s", apiDomain, endpoint)
	// create request
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	// authentication is made via bearer token over http headers.
	tokenHeader := fmt.Sprintf("Bearer %s", apiToken)
	req.Header.Add(APITokenName, tokenHeader)
	req.Header.Add("accept", "application/json")
	req.Header.Add("accept-language", "en-US")
	
	// create client to do the request
	client := http.Client{
		Timeout: time.Second * 45,
	}
	
	// fetch response
	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to user api")
	}
	
	return resp, nil
}
