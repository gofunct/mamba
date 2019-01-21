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

package docker

import (
	"fmt"
	"github.com/gofunct/mamba/docker/commands"
	"github.com/spf13/cobra"
)

var (
	dGraph bool
)

func init() {
	CreateCmd.Flags().BoolVar(&dGraph, "dgraph", false, "start a local dgraph database server")
}

var CreateCmd = createCmd()

func createCmd() *cobra.Command {
	switch {
	case dGraph:
		fmt.Println("running dgraph command...")
		return &cobra.Command{
			Use:   "create",
			Short: "Commands for interacting with a database",
			Run:   commands.RunDgraph(),
		}
	default:
		fmt.Println("running default command...")
		return &cobra.Command{
			Use:   "create",
			Short: "Commands for interacting with a database",
			Run:   commands.DefaultFunc(),
		}
	}

}
