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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// certCmd represents the cert command
var (
	cert    bool
	certCmd = &cobra.Command{
		Use:   "cert",
		Short: "generate a server key and certificate",
		PreRun: func(cmd *cobra.Command, args []string) {
			OsExec(
				"openssl",
				"genrsa",
				"-out",
				"server.key",
				"2048",
			)
		},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("generating certificates...")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			OsExec(
				"openssl",
				"req",
				"-new",
				"-x509",
				"-key",
				"server.key",
				"-out",
				"server.pem",
				"-days",
				"3650",
			)
		},
	}
)

func init() {
	rootCmd.AddCommand(certCmd)
	certCmd.Flags().BoolVar(&cert, "cert", false, "")
}
