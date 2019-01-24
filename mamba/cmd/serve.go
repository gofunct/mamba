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
	"google.golang.org/grpc"
	"net"
)

var (
	port string
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "port  to listen on")
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a grpc server to handle remote script requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}
		var opts []grpc.ServerOption

		grpcServer := grpc.NewServer(opts...)
		//script.RegisterScriptServiceServer(grpcServer, script.NewScriptHandler())
		return grpcServer.Serve(lis)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
