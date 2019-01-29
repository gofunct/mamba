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
	"github.com/gofunct/mamba/runtime/function"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var (
	kConfig string
	home    = homedir.HomeDir()
)

// deployCmd represents the deploy command
var RootCmd = &cobra.Command{
	Use: "contain",
	Short: "🐍 A Docker development utility",
}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "mamba contain pull --image busybox --repo busybox",
	Run:   function.DockerPull(),
}

// containCmd represents the contain command
var imgCmd = &cobra.Command{
	Use: "images",
	Run: function.DockerImages(),
}

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "A brief description of your command",
	Run:   function.RunDgraph(),
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build called")
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(graphCmd)
}
