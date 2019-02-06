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

package walk

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	wConfig, in, out, pkg string
)

// protocCmd represents the protoc command
var RootCmd = &cobra.Command{
	Use:   "walk",
	Short: "🐍 Walk a filepath with a given function and file extension",
}

func init() {
	RootCmd.AddCommand(htmlCmd, goGoCmd, grpcCmd, jsCmd)
	RootCmd.PersistentFlags().StringVar(&wConfig, "config", filepath.Join(os.Getenv("PWD"), "config", "walk.yaml"), "walk config path")
	RootCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "input dir")
	RootCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "output dir")
	RootCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "", "package name")
}
