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

package commands

import (
	"github.com/gofunct/mamba/runtime/scripter/cli"
	mux2 "github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"net/http"
)

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "🐍 One line golang scripts",

	RunE: func(cmd *cobra.Command, args []string) error {
		mux := mux2.NewRouter()
		mux.HandleFunc("/"+cmd.Use, func(writer http.ResponseWriter, request *http.Request) {
			cmd.SetOutput(writer)

			cli.Run(cmd)
		})
		return http.ListenAndServe(":11000", mux)
	},
}

func init() {
	RootCmd.AddCommand(scriptCmd)
}
