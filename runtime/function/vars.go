package function

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/gofunct/common/pkg/logger/zap"
)

var (
	client                                       *docker.Client
	endpoint, containerName, imageName, repoName string
	ports                                        = make(map[docker.Port]struct{})
	shell                                        []string
)

func init() {
	var err error
	client, err = docker.NewClient(endpoint)
	if err != nil {
		zap.LogF("failed to setup docker client", err)

	}
}
