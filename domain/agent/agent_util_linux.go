// +build linux

package agent

import (
	"bytes"
	`encoding/json`
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	
	"github.com/pkg/errors"
	"github.com/usestrix/cli/domain/consts"
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
	
	return errors.New("unknown deployment type. check your agent configuration again")
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
GroupID=root
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
	err = ioutil.WriteFile(filepath.Join(consts.ConfigDir, consts.ConfigFile), data, 0600)
	if err != nil {
		return err
	}
	
	return nil
}

func StopDaemon() error {
	cmd := exec.Command("systemctl", "stop", "strixeyed")
	
	err := cmd.Run()
	if err != nil {
		return err
	}
	
	return nil
}
