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
	kitlog "github.com/go-kit/kit/log"
	"github.com/gofunct/mamba/cmd/gcloud"
	"github.com/gofunct/mamba/cmd/local"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	{
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC, "origin")
		log.SetOutput(kitlog.NewStdlibAdapter(logger))
	}

	{
		rootCmd.AddCommand(gcloud.RootCmd)
		rootCmd.AddCommand(local.RootCmd)
		rootCmd.AddCommand(protocCmd)
		rootCmd.AddCommand(htmlCmd)
		rootCmd.AddCommand(testCmd)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mamba",
	Short: "A general purpose scripting utility for developers and administrators",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
