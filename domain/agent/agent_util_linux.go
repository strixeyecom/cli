//go:build linux
// +build linux

package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	"github.com/strixeyecom/cli/domain/consts"
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
//
// It handles kubernetes/docker based differentiation
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

	return errors.New("unknown deployment type. check your agent configuration again")
}

func (a AgentInformation) CreateServiceFile() error {
	return createServiceFile(a)
}

// createServiceFile creates a ever running service file depending on your configuration.
//
// Different service files are created for Linux/Docker, Linux/Kubernetes or Darwin/Kubernetes.
func createServiceFile(agentInformation AgentInformation) error {
	var (
		err         error
		serviceFile string
	)
	switch agentInformation.Config.Deployment {
	case "docker":
		serviceFile, err = createDockerServiceFile()
	case "kubernetes":
		serviceFile, err = createKubernetesServiceFile()
	}

	// 	save service file
	servicePath := filepath.Join(consts.ServiceDir, consts.ServiceFile)
	err = ioutil.WriteFile(servicePath, []byte(serviceFile), 0600)
	if err != nil {
		return err
	}

	return nil
}

// createServiceFile creates service file, again depending on the environment. On Unix machines,
// StrixEye Daemon uses systemd. On windows machines, who knows?
//
// Although silly, it also controls required tools like mkdir,chown or docker compose and uses given paths.
func createDockerServiceFile() (string, error) {

	var (
		err                                 error
		dockerComposePath, mkdir, chmodPath string
		composeFilePath                     string
	)

	// systemd service file
	serviceFile := `[Unit]
Description=StrixEye Agent Service Daemon

[Service]
User=root
Type=simple
Requires=docker.service

Restart=always
RestartSec=2s

ExecStart=%s
ExecStopPost=%s -f %s down
WorkingDirectory=%s

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true


ExecStartPre=%s -p %s
ExecStartPre=%s 755 %s
StandardError=syslog
SyslogIdentifier=%s

[Install]
WantedBy=multi-user.target
`
	execPath := filepath.Join(consts.DaemonDir, consts.DaemonName)

	// get docker compose path
	dockerComposePath, err = exec.LookPath("docker-compose")
	if err != nil {
		return "", err
	}

	// get chown and chmod path
	chmodPath, err = exec.LookPath("chmod")
	if err != nil {
		return "", err
	}

	mkdir, err = exec.LookPath("mkdir")
	if err != nil {
		return "", err
	}

	composeFilePath = filepath.Join(consts.WorkingDir, consts.DockerComposeFileName)

	return fmt.Sprintf(
		serviceFile,
		execPath, dockerComposePath, composeFilePath, consts.WorkingDir, mkdir,
		consts.LogFile, chmodPath, consts.LogFile, consts.DaemonName,
	), nil
}

func createKubernetesServiceFile() (string, error) {
	var (
		err              error
		mkdir, chmodPath string
	)

	// systemd service file
	serviceFile := `[Unit]
Description=StrixEye Agent Service Daemon

[Service]
User=root
GroupID=root
Type=simple

Restart=always
RestartSec=2s

ExecStart=%s
WorkingDirectory=%s

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true


ExecStartPre=%s -p %s
ExecStartPre=%s 755 %s
StandardError=syslog
SyslogIdentifier=%s

[Install]
WantedBy=multi-user.target
`
	execPath := filepath.Join(consts.DaemonDir, consts.DaemonName)

	// get chown and chmod path
	chmodPath, err = exec.LookPath("chmod")
	if err != nil {
		return "", err
	}

	mkdir, err = exec.LookPath("mkdir")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		serviceFile,
		execPath, consts.WorkingDir, mkdir,
		consts.LogFile, chmodPath, consts.LogFile, consts.DaemonName,
	), nil
}

func InstallCompleted() {
	fmt.Println(
		`Install successfully completed.
Please run the following command immediately

	$ systemctl daemon-reload

Start StrixEye Daemon with
	$ systemctl start strixeyed

To enable start at boot (highly recommended)
	$ systemctl enable strixeyed`,
	)
}

// SaveAgentConfig store agent config in file on *NIX systems.
func SaveAgentConfig(cfg Agent) error {
	// marshal config
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	// save data
	// #nosec
	err = ioutil.WriteFile(filepath.Join(consts.ConfigDir, consts.ConfigFile), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func StopDaemon() error {
	var outbuf, errbuf bytes.Buffer
	exitCode := 0
	stderr := errbuf.String()
	const defaultFailedCode = 1

	cmd := exec.CommandContext(context.Background(), "systemctl", "stop", "strixeyed")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
			if exitCode == 5 {
				return nil
			}
			return err
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			// return errors.Errorf("Could not get exit code for failed program")
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
			return err
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return nil
}

var IsCorrectUser = func() error {
	if !IsRootUser() {
		return errors.New("user not root")
	}
	return nil
}

func IsRootUser() bool {
	return os.Geteuid() == 0
}
