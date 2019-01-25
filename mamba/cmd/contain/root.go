// Copyright © 2019 Coleman Word <coleman.word@gofunct.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package contain

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/gofunct/mamba/pkg/input"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

var (
	endpoint, containerName, imageName, repoName string
	shell                                        []string
	ports                                        = make(map[docker.Port]struct{})
	client                                       *docker.Client
	ui                                           *input.UI

	RootCmd = &cobra.Command{
		Use: "contain",
	}
)

// Execute executes the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute:%s/n", err)
	}
}

func init() {
	var err error
	cobra.OnInitialize(
		func() {
			client, err = docker.NewClient(endpoint)
			if err != nil {
				log.Fatalln("failed to setup docker client", err.Error())

			}
		},
		func() {
			ui = input.DefaultUI()
		},
	)

	{
		RootCmd.AddCommand(imgCmd)
		RootCmd.AddCommand(bindCmd)
		RootCmd.AddCommand(buildCmd)
		RootCmd.AddCommand(tlsCmd)
		RootCmd.AddCommand(execCmd)
		RootCmd.AddCommand(pushCmd)
		RootCmd.AddCommand(pullCmd)
		RootCmd.AddCommand(pullCmd)
		RootCmd.AddCommand(storeCmd)
		RootCmd.AddCommand(inspectCmd)
	}
	{
		RootCmd.PersistentFlags().StringSliceVar(&shell, "shell", []string{"echo", "hello world"}, "inter-container commands")
		RootCmd.PersistentFlags().StringVar(&repoName, "repo", "", "repo name")
		RootCmd.PersistentFlags().StringVar(&imageName, "image", "", "image name")
		RootCmd.PersistentFlags().StringVar(&containerName, "container", "", "container name")
		RootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "endpoint for docker client to connect to")
	}
}
