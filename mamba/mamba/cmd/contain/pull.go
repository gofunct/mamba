// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package contain

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"log"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "mamba contain pull --image busybox --repo busybox",
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.PullImage(docker.PullImageOptions{Repository: repoName}, docker.AuthConfiguration{}); err != nil {
			log.Fatal(err.Error())
		}
		ports["8080/tcp"] = struct{}{}
		config := docker.Config{
			ExposedPorts: ports,
			Env:          nil,
			Cmd:          shell,
			Image:        imageName,
		}
		opts := docker.CreateContainerOptions{Name: containerName, Config: &config}
		if cont, err := client.CreateContainer(opts); err != nil {
			log.Fatalln(err.Error())
		} else {
			if err := client.StartContainer(cont.ID, &docker.HostConfig{}); err != nil {
				log.Fatalln(err.Error())
			}
		}

	},
}
