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
	"fmt"
	"github.com/gofunct/mamba/docker"
	"github.com/gofunct/mamba/generator"
	"github.com/gofunct/mamba/manager/logging"
	"github.com/gofunct/mamba/static"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	in  string
	out string
	pkg string
)

func init() {
	{
		logging.AddLoggingFlags(rootCmd)
		rootCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "path to input directory")
		rootCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "path to output directory")
		rootCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "", "package name")

	}

	{
		rootCmd.AddCommand(generator.GoGoCmd)
		rootCmd.AddCommand(static.RootCmd)
		rootCmd.AddCommand(testCmd)
		rootCmd.AddCommand(generator.ProtoGenCmd)
		rootCmd.AddCommand(serveCmd)
		rootCmd.AddCommand(docker.RootCmd)
		rootCmd.AddCommand(inputCmd)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mamba",
	Short: "A general purpose scripting utility for developers and administrators",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s\n", errors.WithStack(err))
	}
}
