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
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"strings"
	"time"
)

func init() {
	logger = logging.NewLogCtx(logrus.New())
	query = input.DefaultUI()
	for _, x := range initializers {
		x()
	}
}

var (
	logger       *logging.CtxLogger
	query        *input.UI
	initializers []func()
)

const EnablePrefixMatching = false
const EnableCommandSorting = true
const MousetrapHelpText string = `This is a command line tool.

You need to open cmd.exe and run it from there.
`

type Command struct {
	Version      string
	Dependencies []string
	PreRun       MambaFunc
	Login        http.HandlerFunc
	Home         http.HandlerFunc
	FAQ          http.HandlerFunc
}

func (c *Command) Execute(ctx context.Context) error {
	c.AddLogging()
	if len(c.Dependencies) > 0 {
		c.SyncRequirements()
	}

	if c.PreRun != nil {
		c.PreRun(c, ctx)
	}
	if c.Login != nil {
		router.Handle("/login", c.Login)
	}
	if c.Home != nil {
		router.Handle("/", c.Home)
	}
	if c.FAQ != nil {
		router.Handle("/faq", c.FAQ)
	}
	router.HandleFunc("/settings", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, fmt.Sprintf("%#v", c.GetMeta()))
	})
	router.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
	router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	var srv *http.Server
	if router == nil {
		router = mux.NewRouter()
	}

	if port := os.Getenv("MAMBA_PORT"); port != "" {
		srv = &http.Server{
			Handler: router,
			Addr:    port,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	} else {
		srv = &http.Server{
			Handler: router,
			Addr:    "127.0.0.1:8080",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	}
	fmt.Println("🐍 starting server on:", srv.Addr)
	fmt.Println("🐍 type Ctrl-C to shutdown ", srv.Addr)

	return srv.ListenAndServe()
}

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
