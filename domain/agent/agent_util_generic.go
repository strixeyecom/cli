package agent

import (
	`context`
	
	`github.com/docker/docker/api/types`
	"github.com/docker/docker/client"
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	generic controllers that work over package apis.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// checkDockerInstalled uses docker api, it is os independent to us. IDK what docker does under the hood.
func checkDockerRunning() error {
	var (
		err error
		ctx = context.Background()
	)
	
	// create a client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	
	// try to list containers
	_, err = cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	
	return nil
}
