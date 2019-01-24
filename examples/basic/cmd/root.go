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
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func init() {
	ctx = context.TODO()
}

var (
	ctx context.Context
)

var root = &mamba.Command{
	Version:      "v0.1.1",
	Dependencies: nil,
	PreRun: func(svc *mamba.Command, ctx context.Context) {
		fmt.Println("Welcome " + os.Getenv("USER") + "!")
	},
	Hanldlers: map[string]http.HandlerFunc{
		"": func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "this is where your web app will be located")
		},
		"login": func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "this is where your users will login")
		},
		"faq": func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "this is where your users will go for help")
		},
	},
}

func Execute() {
	if err := root.Execute(ctx); err != nil {
		fmt.Printf("%#v", errors.WithStack(err))
	}
}
