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

package server

import (
	"context"
	"fmt"
	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
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
	for _, x := range initializers {
		x()
	}
}

var (
	g            group.Group
	initializers []func()
)

const EnablePrefixMatching = false
const EnableServiceSorting = true
const MousetrapHelpText string = `This is a command line tool.

You need to open cmd.exe and run it from there.
`

type Service struct {
	Version      string
	Dependencies []string
	Hanldlers    map[string]http.HandlerFunc
	Options      []grpc.ServerOption
	PostRun      MambaFunc
	r            *mux.Router
}

func (c *Service) Execute(ctx context.Context) error {

	grpcServer := grpc.NewServer(c.Options...)
	for k, v := range c.Hanldlers {
		c.r.Handle("/"+k, v)
	}

	c.r.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
	c.r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	c.r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	c.r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	c.r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	c.r.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	var srv *http.Server
	if c.r == nil {
		c.r = mux.NewRouter()
	}

	if port := os.Getenv("MAMBA_PORT"); port != "" {
		srv = &http.Server{
			Handler: c.handleGrpc(grpcServer),
			Addr:    port,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	} else {
		srv = &http.Server{
			Handler: c.handleGrpc(grpcServer),
			Addr:    "127.0.0.1:8080",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	}
	fmt.Println("🐍 starting server on:", srv.Addr)
	fmt.Println("🐍 type Ctrl-C to shutdown ", srv.Addr)
	g.Add(func() error {
		zap.Debug("starting server", "transport", "server/HTTP")
		zap.Debug("starting server", "addr", srv.Addr)
		return srv.ListenAndServe()
	}, func(error) {
		srv.Shutdown(ctx)
	})
	if c.PostRun != nil {
		srv.RegisterOnShutdown(func() {
			c.PostRun(c, ctx)
		})
	}
	return g.Run()
}

type MambaFunc func(command *Service, ctx context.Context)

func OnInitialize(y ...func()) {
	initializers = append(initializers, y...)
}

func (m *Service) WriteFile(f string, d []byte) error {
	return ioutil.WriteFile(f, d, 0755)
}

func (m *Service) ReadFile(f string) ([]byte, error) {
	return ioutil.ReadFile(f)
}

func (m *Service) ReadStdIn() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}
func (m *Service) ReadReader(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func (m *Service) ReadDir(f string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(f)
}

func (c *Service) ExecString(args ...string) (stdout string, err error) {
	stdoutb, err := c.ExecBytes(args...)
	return strings.TrimSpace(string(stdoutb)), err
}

func (c *Service) ExecBytes(args ...string) (stdout []byte, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return stdoutb, nil
}

func (c *Service) OsExec(args ...string) {
	s, err := c.ExecString(args...)
	if s != "" {
		if _, err := fmt.Fprintf(os.Stderr, s); err != nil {
			zap.LogF("write output to stderr", err)
		}
	}
	zap.LogF("os exec command", err)

}
