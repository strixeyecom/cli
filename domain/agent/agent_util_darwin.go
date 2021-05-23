// +build darwin

package agent

import (
	`bytes`
	`fmt`
	`os/exec`
	
	`github.com/pkg/errors`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	utility functions for StrixEye agents running on Linux.
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// CheckIfHostSupports controls whether you can install your current agent on the host machine or not.
func (a AgentInformation) CheckIfHostSupports() error {
	var (
		err error
	)
	
	// Run generic checks
	
	// 	control according to deployment type
	// check for docker/docker-compose deployment.
	if a.Config.Deployment == "docker" {
		// check whether docker compose is installed.
		err = checkIfDockerComposeExists()
		if err != nil {
			return errors.Wrap(
				err, "can not find docker compose, "+
					"are you sure you have installed both docker and its tool docker compose?",
			)
		}
		
		// check whether there is a running docker daemon.
		err = checkDockerRunning()
		if err != nil {
			return errors.Wrap(
				err, "docker is not running or maybe cli missing permissions. "+
					"Check your docker configuration",
			)
		}
		
		return nil
	}
	
	// check for kubectl api
	if a.Config.Deployment == "kubernetes" {
		return nil
	}
	
	return errors.New("unknown deployment type. check your configuration again")
}

func checkIfDockerComposeExists() error {
	cmd := exec.Command("docker-compose", "version")
	
	var output bytes.Buffer
	cmd.Stdout = &output
	
	err := cmd.Run()
	// if exit code != 0, it means docker-compose not found.
	if err != nil {
		return err
	}
	
	fmt.Println(output.String())
	return nil
}
