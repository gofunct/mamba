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
	"github.com/gofunct/mamba/manager/input"
	manager2 "github.com/gofunct/mamba/manager"
	"github.com/spf13/cobra"
)

var UI = &input.UI{
	Queries: []*input.Query{
		{
		Q: "What is your favorite restaurant?",
		Opts: &input.Options{
			ValidateFunc: func(s string) error {
				if s == "" {
					return input.ErrEmpty
				}
				if len(s) > 50 {
					return input.ErrOutOfRange
				}
				return nil
			},
			Default: "Mcdonalds",
			Required: true,
			Loop: true,
		},
		Tag: "restaurant",
		},
	},
}

// inputCmd represents the input command
var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "temporary",

	Run: func(cmd *cobra.Command, args []string) {
		var manager = manager2.NewManager([]string{"name", "email"})
		manager.AddFlagSet(cmd.Flags())
		manager.Q.AddQueries(UI.Queries...)
		manager.Q.Query()
		manager.SyncRequirements()
		manager.Debug()
	},
}