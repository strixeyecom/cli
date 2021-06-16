// +build !develop,windows

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

	ConfigDir = "/etc/strixeye/config"
	ConfigFile = "config.json"

	DownloadHost = "https://downloads.strixeye.com"
	APIHost = "https://api.strixeye.com"

	DockerComposeFileName = "docker-compose.yml"

	LogFile = "/var/log/strixeyed"

	ServiceDir  = "/lib/systemd/system"
	ServiceFile = "strixeyed.service"
	
	DownloadZipName = "manager.tar.gz"

)

// global variables (not cool) for this file
var ()
