package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	ssesvr "github.com/mt5225/ct-runner/ssesvr"
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

var keepAliveCMD = []string{"sleep", "infinity"}

// Run container
func (cmd *Command) Run(r *gin.Engine) (string, error) {
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
			Image: cmd.Image,
			Cmd:   keepAliveCMD,
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

	fmt.Println("Create container ID = ", cont.ID)

	// start container
	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	if err != nil {
		panic(err)
	}

	// run commands on container
	resp, err := cli.ContainerExecCreate(ctx, cont.ID, types.ExecConfig{
		Cmd:          cmd.Commands,
		Env:          cmd.Env,
		AttachStdout: true,
		AttachStderr: true,
	})

	if err != nil {
		panic(err)
	}

	resultStream, err := inspectExecResp(ctx, resp.ID)

	//output the web page
	ssesvr.SseServer(r, resultStream)

	defer cli.ContainerKill(ctx, cont.ID, "SIGKILL")

	return fmt.Sprint(cont.ID), err
}

// TODO streaming to websocket
func inspectExecResp(ctx context.Context, id string) (*types.HijackedResponse, error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	defer docker.Close()

	resp, err := docker.ContainerExecAttach(ctx, id, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	return &resp, nil
}

type execResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}
