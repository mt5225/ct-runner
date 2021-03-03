package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	guuid "github.com/google/uuid"
)

//Command to run docker image
type Command struct {
	Image    string
	Env      []string
	Commands []string
	Cached   []string
	Shell    string
	Workdir  string
	Stash    map[string]string
	UnStash  map[string]string
}

// Run container
func (cmd *Command) Run() (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	id := guuid.New()
	ctx := context.Background()

	// create container
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        cmd.Image,
			Cmd:          cmd.Commands,
			Env:          cmd.Env,
			AttachStdout: true,
			AttachStderr: true,
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		nil,
		nil,
		id.String())

	if err != nil {
		panic(err)
	}

	// start container
	if err := cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("Container is started", cont.ID)

	// defer cli.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{Force: true})

	options := types.ContainerLogsOptions{ShowStdout: true}

	out, err := cli.ContainerLogs(ctx, cont.ID, options)
	if err != nil {
		panic(err)
	}

	// defer out.Close()

	io.Copy(os.Stdout, out)
	return cont.ID, nil
}
