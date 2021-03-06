package runner

import (
	"context"
)

// Command to run docker image
type Command struct {
	Image    string
	Env      []string
	Commands []string
	Cached   []string
	Shell    string
	Workdir  string
	Stash    map[string]string
	UnStash  map[string]string

	// for container
	ContainerInstance *Container
}

// Container Info
type Container struct {
	Context context.Context
	ID      string
	RunID   string
}
