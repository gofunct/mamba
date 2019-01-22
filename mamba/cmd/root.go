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
	"github.com/gofunct/mamba"
	"github.com/pkg/errors"
)

var (
	userLicense string
	in          string
	out         string
	pkg         string
	rootCmd     = &mamba.Command{
		Use:     "mamba",
		Version: "",
		Aliases: nil,
		Info: `Mamba is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Mamba application.`,
		Hidden: false,
		Env:    nil,
		// Args set in ValidArgs will be set via query if not found
		ValidArgs: nil,
		Args:      nil,
		UsageF:    nil,
		UsageTmpl: "",
		PreRun:    nil,
		Run:       nil,
		PostRun:   nil,
		// Useful for passing args to os.Exec
		DisableFlagParsing: false,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Fatalf("%s\n", errors.WithStack(err))
	}
}

func init() {
	{
		rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
		rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	}

	rootCmd.Set("author", "NAME HERE <EMAIL ADDRESS>")
	rootCmd.Set("license", "apache")

	{
		rootCmd.AddCommand(walkCmd)
		rootCmd.AddCommand(scriptCmd)
		rootCmd.AddCommand(addCmd)
		rootCmd.AddCommand(initCmd)
		rootCmd.AddCommand(testCmd)
	}

}
