package common

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io/ioutil"
	"log"
	"strings"
)

func UpdateGradingImage() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	CheckIfError(err)

	// Update docker image
	_, err = cli.ImagePull(ctx, "cliftonzac/autograder:latest", types.ImagePullOptions{})
	CheckIfErrorWithMessage(err, "Could not update image.")
}

func RemoveContainer(client *client.Client, containerName string) error {
	ctx := context.Background()

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := client.ContainerRemove(ctx, containerName, removeOptions); err != nil {
		log.Printf("Unable to remove container: %s", err)
		return err
	}

	return nil
}

// StopAndRemoveContainer Stop and remove a docker container
func StopAndRemoveContainer(client *client.Client, containerName string) error {
	ctx := context.Background()

	if err := client.ContainerStop(ctx, containerName, nil); err != nil {
		log.Printf("Unable to stop container %s: %s", containerName, err)
	}

	return RemoveContainer(client, containerName)

}

func RunDockerCommand(client *client.Client, containerName string, cmdToExecute []string, workingDir string, print bool) {
	ctx := context.Background()

	respID, err := client.ContainerExecCreate(ctx, containerName, types.ExecConfig{
		Detach:       false,
		Privileged:   false,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
		WorkingDir:   workingDir,
		Cmd:          cmdToExecute,
	})

	CheckIfError(err)

	response, err := client.ContainerExecAttach(context.Background(), respID.ID, types.ExecStartCheck{})
	CheckIfError(err)

	defer response.Close()

	//for {
	//	inspect, err := client.ContainerExecInspect(context.Background(), respID.ID)
	//	CheckIfError(err)
	//
	//	if !inspect.Running {
	//		break
	//	}
	//
	//}

	// DO NOT REMOVE; This ensures the task is completed before moving on to the next command.
	data, _ := ioutil.ReadAll(response.Reader)

	if print {
		fmt.Println(string(data))
	}
}

func ForceStartContainer(client *client.Client, containerName string, print bool) {

	ctx := context.Background()

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{})
	CheckIfError(err)

	for _, indContainer := range containers {

		for _, name := range indContainer.Names {
			if strings.Trim(name, "/") == containerName {

				if indContainer.Status == "running" {
					err := StopAndRemoveContainer(client, containerName)
					CheckIfError(err)
				} else {
					err := RemoveContainer(client, containerName)
					CheckIfError(err)
				}
				break
			}
		}
	}

	// Create Container
	resp, err := client.ContainerCreate(ctx,
		&container.Config{
			Image: "cliftonzac/autograder:latest",
			Cmd:   []string{"bash"},
			Tty:   true,
		},
		nil,
		nil,
		nil,
		containerName,
	)
	CheckIfError(err)

	err = client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	CheckIfError(err)

}
