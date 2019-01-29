package function

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/gofunct/common/pkg/encode"
	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/gofunct/mamba/runtime"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"net/http"
)

func DockerPull() runtime.Function {
	return func(cmd *cobra.Command, args []string) {
		zap.LogF("pulling image", client.PullImage(docker.PullImageOptions{Repository: repoName}, docker.AuthConfiguration{}))
		ports["8080/tcp"] = struct{}{}
		config := docker.Config{
			ExposedPorts: ports,
			Env:          nil,
			Cmd:          shell,
			Image:        imageName,
		}
		opts := docker.CreateContainerOptions{Name: containerName, Config: &config}
		if cont, err := client.CreateContainer(opts); err != nil {
			zap.LogF("create container", err)
		} else {
			if err := client.StartContainer(cont.ID, &docker.HostConfig{}); err != nil {
				zap.LogF("start container", err)
			}
		}

	}
}

func DockerImages() runtime.Function {
	return func(cmd *cobra.Command, args []string) {
		imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
		if err != nil {
			panic(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/images", func(writer http.ResponseWriter, request *http.Request) {
			for i, img := range imgs {

				cmd.SetOutput(writer)
				cmd.Println(i, encode.PrettyJsonString(img))
			}
		})

		log.Info("Starting contain server on:", "0.0.0.0:10000")
		log.Info("Press Crtl-C to shutdown...")
		log.Infof("Routes:/n%s", "/images")
		if err := http.ListenAndServe("0.0.0.0:10000", router); err != nil {
			log.Fatalf("%s\n", err)
		}

	}
}
