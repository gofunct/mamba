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
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"

	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		endpoint := "unix:///var/run/docker.sock"
		client, _ := docker.NewClient(endpoint)
		container, err := client.InspectContainer(containerName)
		if err != nil {
			fmt.Printf("failed to inspect container %v", err)
			os.Exit(1)
		}
		state := container.State
		fmt.Println("ID: ", container.ID)
		fmt.Println("Pid: ", state.Pid)
	},
}

func init() {

}
