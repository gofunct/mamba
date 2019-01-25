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
	"github.com/gofunct/mamba/pkg/encode"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"net/http"
)

func init() {
}

// containCmd represents the contain command
var imgCmd = &cobra.Command{
	Use: "images",
	Run: func(cmd *cobra.Command, args []string) {
		imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
		if err != nil {
			panic(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/images", func(writer http.ResponseWriter, request *http.Request) {
			for i, img := range imgs {

				cmd.SetOutput(writer)
				cmd.Println(i, encode.PrettyJson(img))
			}
		})

		log.Info("Starting contain server on:", "0.0.0.0:10000")
		log.Info("Press Crtl-C to shutdown...")
		log.Infof("Routes:/n%s", "/images")
		if err := http.ListenAndServe("0.0.0.0:10000", router); err != nil {
			log.Fatalf("%s\n", err)
		}

	},
}
