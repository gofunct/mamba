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

package cmd

import (
	"context"
	"fmt"
	"github.com/gofunct/mamba/api/script"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	cmdName string
	cmdDir  string
	cmdArgs []string
	host    string
)

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&cmdName, "name", "n", "", "")
	execCmd.Flags().StringVarP(&cmdDir, "dir", "d", "", "")
	execCmd.Flags().StringSliceVarP(&cmdArgs, "args", "a", []string{}, "")
	execCmd.Flags().StringVarP(&host, "host", "h", "localhost:8080", "")
}

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "execute a script on a backend grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		// Create a insecure gRPC channel to communicate with the server.
		conn, err := grpc.Dial(
			host,
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		client := script.NewScriptServiceClient(conn)
		out, err := client.Exec(context.Background(), &script.Command{
			Name: cmdName,
			Dir:  cmdDir,
			Args: cmdArgs,
			Env:  os.Environ(),
		})
		if err != nil {
			fmt.Println(string(out.Out))
			log.Fatalln("failed to execute command", errors.WithStack(err))
		}
		fmt.Println(string(out.Out))
	},
}
