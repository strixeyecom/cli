package agent

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

/*
	Created by aomerk at 8/9/21 for project strixeye
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func RemoveDockerVolumeByName(name string) error {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	vList, err := dockerClient.VolumeList(ctx, filters.Args{})
	if err != nil {
		return err
	}
	for _, volume := range vList.Volumes {
		if volume.Name != name {
			continue
		}

		// remove given docker volume
		err = dockerClient.VolumeRemove(ctx, volume.Name, true)
		if err != nil {
			return err
		}
	}

	return nil
}
