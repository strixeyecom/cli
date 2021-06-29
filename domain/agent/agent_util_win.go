// +build windows

package agent

import (
	`github.com/pkg/errors`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	utility functions for StrixEye agents running on Linux.
*/


// checkIfAnotherAgentRunning tries to find a running strixeyed daemon and returns nil if **no agent is
// running**
//
// Control mechanism depends on system, but in general, to avoid false positives,
// we use a dedicated PID file to keep track of a running StrixEye daemon.
func checkIfAnotherAgentRunning() error {
	// A StrixEye Daemon creates a pid file to show that it is running.
	//
	// We should check if such file exists.
	// There are cases where strixeyed doesn't shut down gracefully and leave a strixeyed pid behind
	_, err := os.Stat(consts.PidFile)
	if err == nil {
		return ErrAnotherAgentRunning
	}
	
	// If the error is a file not found/not exists error,
	// it means that there are no strixeyed running on host machine.
	
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	
	// no other strixeyed alive
	return nil
}


// checkIfHostSupports controls whether you can install your current agent on the host machine or not.
func (a AgentInformation) checkIfHostSupports() error {
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