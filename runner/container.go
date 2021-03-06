package runner

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
)

var keepAliveCMD = []string{"sleep", "infinity"}

// Create container
func (cmd *Command) Create() error {
	cli, err := client.NewEnvClient()

	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	id := guuid.New()

	// create container
	cont, err := cli.ContainerCreate(
		cmd.ContainerInstance.Context,
		&container.Config{
			Image: cmd.Image,
			Cmd:   keepAliveCMD,
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		nil,
		nil,
		id.String())

	cmd.ContainerInstance.ID = cont.ID
	fmt.Println("Create container ID = ", cont.ID)
	return err
}

// Run container
func (cmd *Command) Run(r *gin.Engine) error {
	cli, err := client.NewEnvClient()

	// start container
	err = cli.ContainerStart(context.Background(), cmd.ContainerInstance.ID, types.ContainerStartOptions{})

	if err != nil {
		panic(err)
	}

	// run commands on container
	resp, err := cli.ContainerExecCreate(cmd.ContainerInstance.Context, cmd.ContainerInstance.ID, types.ExecConfig{
		Cmd:          cmd.Commands,
		Env:          cmd.Env,
		AttachStdout: true,
		AttachStderr: true,
	})

	if err != nil {
		panic(err)
	}

	cmd.ContainerInstance.RunID = resp.ID
	return err
}
