// +build !develop,darwin

package consts

/*
	Created by aomerk at 6/14/21 for project cli
*/

/*
	strixeyed/strixeye related constant variables
*/

// global constants for file
const (
	DaemonDir  = "/usr/bin"
	DaemonName = "strixeyed"

	WorkingDir = "/etc/strixeye"

	ConfigDir  = "/etc/strixeye/config"
	ConfigFile = "config.json"

	DownloadHost = "downloads.strixeye.com"
	DockerRegistry = "docker.strixeye.com"
	APIHost      = "api.strixeye.com"

	DockerComposeFileName = "docker-compose.yml"

	LogFile = "/var/log/strixeyed"
	PidFile = "/var/run/strixeyed.pid"

	ServiceDir  = "/lib/systemd/system"
	ServiceFile = "strixeyed.service"

	DownloadZipName = "manager.tar.gz"
)

// global variables (not cool) for this file
var ()
