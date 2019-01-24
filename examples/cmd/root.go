// Copyright © 2019 coleman
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
	"context"
	"fmt"
	"github.com/gofunct/mamba"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net/http"
)

func init() {
	ctx = context.TODO()
}

var (
	ctx context.Context
)

var root = &mamba.Command{
	Version: "v0.1.1",
	// You will be prompted to set these environmental variables if they are not found
	Dependencies: []string{"MAMBA_PORT"},
	// Scripts run as they would from the terminal. They run before the server starts
	Scripts: [][]string{
		// This is just an example
		[]string{"echo", "vendoring dependencies..."},
		[]string{"go", "mod", "vendor"},
	},
	// a map of a handler path(without a "/") and a handlerfunc
	// these handlers are served after the scripts finish successfully
	Hanldlers: map[string]http.HandlerFunc{
		"": func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "this is the base url")
		},
		"login": func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "this is where your users could login")
		},
		"faq": func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "this is where your users could go for help")
		},
	},
	// the grpc server handles all grpc requests on the same port as the http handlers
	// this is where you specify grpc middleware
	Options: []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(),
		),
	},
	// PostRun runs after the server has shutdown successfully
	PostRun: func(svc *mamba.Command, ctx context.Context) {
		fmt.Println("server shutdown successfully!")
	},
}

func Execute() {
	// func main() calls this function to execute the root command
	if err := root.Execute(ctx); err != nil {
		fmt.Printf("%#v", errors.WithStack(err))
	}
}
