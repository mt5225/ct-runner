package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

//Command to run docker image
type Command struct {
	Image    string
	Env      map[string]string
	Commands []string
	Cached   []string
	Shell    string
	Workdir  string
	Stash    map[string]string
	UnStash  map[string]string
}

// TerraformImage url
const TerraformImage = "radut/terraform-ansible:latest"

// Run container
func (cmd Command) Run() (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: TerraformImage,
		},
		nil,
		nil,
		nil,
		"UUID")
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", cont.ID)
	return cont.ID, nil
}
