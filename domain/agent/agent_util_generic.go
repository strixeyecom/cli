package agent

import (
	`bytes`
	`context`
	`os/exec`
	
	`github.com/docker/docker/api/types`
	"github.com/docker/docker/client"
	`github.com/fatih/color`
	`github.com/pkg/errors`
	`github.com/usestrix/cli/domain/consts`
)

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	generic controllers that work over package apis.
*/

// global constants for file
const (
	// LogCheckPID verbose message to say checking for running strixeyed
	LogCheckPID = "Controlling strixeyed pid file at " + consts.PidFile
	
	logCheckSupport = "Verifying that host machine supports running StrixEye Agent depending on your agent" +
		" configuration"
	
	logVerificationFailed     = "\t❌❌❌Verification failed. ❌❌❌"
	logVerificationSuccessful = "\tVerification successfully completed. ✅"
	
	DockerDatabaseVolumeName = "strixeye_strixeye-database"
	DockerBrokerVolumeName   = "strixeye_strixeye-queue"
)

var (
	// ErrAnotherAgentRunning tells that there is a strixeyed process somewhere in host machine.
	ErrAnotherAgentRunning = errors.New("another strixeyed is still running")
)

// CheckIfAnotherAgentRunning tries to find a running strixeyed daemon and returns nil if **no agent is
// running**
//
// Control mechanism depends on system, but in general, to avoid false positives,
// we use a dedicated PID file to keep track of a running StrixEye daemon.
func CheckIfAnotherAgentRunning() error {
	color.Blue("Verifying there are no active StrixEye Agents on host machine")
	err := checkIfAnotherAgentRunning()
	if err != nil {
		color.Red("\tVerification failed.")
		return err
	}
	
	color.Yellow("\tVerification successful")
	return nil
}

// CheckIfHostSupports controls whether you can install your current agent on the host machine or not.
func CheckIfHostSupports(a AgentInformation) error {
	return a.CheckIfHostSupports()
}

// CheckIfHostSupports controls whether you can install your current agent on the host machine or not.
func (a AgentInformation) CheckIfHostSupports() error {
	color.Blue(logCheckSupport)
	
	// 	following code depends on the host machine setup (os/arch)
	err := a.checkIfHostSupports()
	if err != nil {
		color.Red(logVerificationFailed)
		return err
	}
	
	color.Yellow(logVerificationSuccessful)
	return nil
}

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

func checkIfDockerComposeExists() error {
	cmd := exec.Command("docker-compose", "version")
	
	var output bytes.Buffer
	cmd.Stdout = &output
	
	err := cmd.Run()
	
	// if exit code != 0, it means docker-compose not found.
	if err != nil {
		return err
	}
	
	color.Green(output.String())
	return nil
}
