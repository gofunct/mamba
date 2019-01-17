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
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var md = make(map[string]string)
var tmplDir string

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:         "replace",
	Short:       "for all files with the .tmpl extension, replace the keys with the values present in the provided metadata",
	Annotations: md,
	Run: func(cmd *cobra.Command, args []string) {
		if err := replaceFunc(tmplDir); err != nil {
			log.Fatalln("failed to run command", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)
	replaceCmd.PersistentFlags().StringToStringVarP(&md, "meta", "m", nil, "key value pairs to add as metadata")
	replaceCmd.PersistentFlags().StringVarP(&tmplDir, "dir", "d", "", "template directory")
}

func replaceFunc(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".tmpl" {
			bytes, err := ioutil.ReadFile(info.Name())
			if err != nil {
				return errors.WithStack(err)
			}
			stringer := string(bytes)

			for old, new := range md {
				stringer = strings.Replace(stringer, old, new, -1)
			}
			if err := ioutil.WriteFile(info.Name(), []byte(stringer), 0666); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}
