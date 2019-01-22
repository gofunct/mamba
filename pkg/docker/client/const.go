package client

import "github.com/fsouza/go-dockerclient"

const endpoint = "unix:///var/run/docker.sock"

var (
	Client, _ = docker.NewClient(endpoint)
)
