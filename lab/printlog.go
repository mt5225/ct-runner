package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true}

	out, err := cli.ContainerLogs(ctx, "3249740e9b50364e2051fcd18789a2b0b32a0d484b79378873af04ba36dcbbce", options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
