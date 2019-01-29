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

package ctl

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var (
	kConfig string
	home    = homedir.HomeDir()
)

// deployCmd represents the deploy command
var RootCmd = &cobra.Command{
	Use: "ctl",
	Short: "🐍 A Kubernetes development utility",
}

func init() {
	RootCmd.AddCommand(deployCmd)
	RootCmd.PersistentFlags().StringVar(&kConfig, "config", filepath.Join(home, ".kube", "config"), "kube config path")
}
