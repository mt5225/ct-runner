package container

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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

var keepAliveCMD = []string{"sleep", "infinity"}

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

	result, err := inspectExecResp(ctx, resp.ID)

	defer cli.ContainerKill(ctx, cont.ID, "SIGKILL")

	return fmt.Sprint(result), err
}

// TODO streaming to websocket
func inspectExecResp(ctx context.Context, id string) (execResult, error) {
	var execResult execResult
	docker, err := client.NewEnvClient()
	if err != nil {
		return execResult, err
	}
	defer docker.Close()

	resp, err := docker.ContainerExecAttach(ctx, id, types.ExecStartCheck{})
	if err != nil {
		return execResult, err
	}
	defer resp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return execResult, err
		}
		break

	case <-ctx.Done():
		return execResult, ctx.Err()
	}

	stdout, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		return execResult, err
	}
	stderr, err := ioutil.ReadAll(&errBuf)
	if err != nil {
		return execResult, err
	}

	res, err := docker.ContainerExecInspect(ctx, id)
	if err != nil {
		return execResult, err
	}

	execResult.ExitCode = res.ExitCode
	execResult.StdOut = string(stdout)
	execResult.StdErr = string(stderr)
	return execResult, nil
}

type execResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}
