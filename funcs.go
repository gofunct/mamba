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

package mamba

import (
	"fmt"
	"github.com/gofunct/mamba/pkg/input"
	"github.com/gofunct/mamba/pkg/logging"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"unicode"
)

func inti() {
	logger = logging.NewLogCtx(logrus.New())
	query = input.DefaultUI()
}

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"rpad":                    rpad,
}

var initializers []func()
var logger *logging.CtxLogger
var query *input.UI

const EnablePrefixMatching = false
const EnableCommandSorting = true
const MousetrapHelpText string = `This is a command line tool.

You need to open cmd.exe and run it from there.
`

type MambaFunc func(command *Command, args []string)

func AddTemplateFunc(name string, tmplFunc interface{}) {
	templateFuncs[name] = tmplFunc
}

func AddTemplateFuncs(tmplFuncs template.FuncMap) {
	for k, v := range tmplFuncs {
		templateFuncs[k] = v
	}
}

func OnInitialize(y ...func()) {
	initializers = append(initializers, y...)
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

// ld compares two strings and returns the levenshtein distance between them.
func ld(s, t string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (m *Command) WriteFile(f string, d []byte) error {
	return ioutil.WriteFile(f, d, 0755)
}

func (m *Command) ReadFile(f string) ([]byte, error) {
	return ioutil.ReadFile(f)
}

func (m *Command) ReadStdIn() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}
func (m *Command) ReadReader(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func (m *Command) ReadDir(f string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(f)
}

func (c *Command) ExecString(args ...string) (stdout string, err error) {
	stdoutb, err := c.ExecBytes(args...)
	return strings.TrimSpace(string(stdoutb)), err
}

func (c *Command) ExecBytes(args ...string) (stdout []byte, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return stdoutb, nil
}

func (c *Command) OsExec(args ...string) {
	s, err := c.ExecString(args...)
	if s != "" {
		if _, err := fmt.Fprintf(os.Stderr, s); err != nil {
			c.Warnf("%s\n%s", "failed to write output to stderr", err)
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
