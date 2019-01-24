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
	"context"
	"fmt"
	"github.com/gofunct/mamba/pkg/input"
	"github.com/gofunct/mamba/pkg/logging"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func inti() {
	logger = logging.NewLogCtx(logrus.New())
	query = input.DefaultUI()

}

var logger *logging.CtxLogger
var query *input.UI

const EnablePrefixMatching = false
const EnableCommandSorting = true
const MousetrapHelpText string = `This is a command line tool.

You need to open cmd.exe and run it from there.
`

type MambaFunc func(command *Command, ctx context.Context)

func OnInitialize(y ...func()) {
	initializers = append(initializers, y...)
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
