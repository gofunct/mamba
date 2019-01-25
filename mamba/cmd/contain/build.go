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
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gofunct/mamba/pkg/logging"
	"github.com/pkg/errors"
	"time"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use: "build",
	Run: BuildImage(),
}

func BuildImage() func(command *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		t := time.Now()
		inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
		tr := tar.NewWriter(inputbuf)
		err := tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 59, ModTime: t, AccessTime: t, ChangeTime: t})
		if err != nil {
			logging.L.Fatalln("failed to write header to dockerfile", errors.WithStack(err).Error())
		}
		i, err := tr.Write([]byte(fmt.Sprintf("FROM alpine \n%s\n%s\n", "COPY . .", `ENTRYPOINT [ "echo", "hello world" ]`)))
		if err != nil {
			logging.L.Fatalln("failed to write header to dockerfile", errors.WithStack(err).Error(), i)
		}

		err = tr.Close()
		if err != nil {
			logging.L.Fatalln("failed to write header to dockerfile", errors.WithStack(err).Error())
		}
		opts := docker.BuildImageOptions{
			Name:         "test2",
			InputStream:  inputbuf,
			OutputStream: outputbuf,
		}
		if err := client.BuildImage(opts); err != nil {
			logging.L.Fatalln("failed to build image", errors.WithStack(err).Error())
		}
	}
}
