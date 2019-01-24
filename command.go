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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

func init() {
	for _, x := range initializers {
		x()
	}
}

var (
	initializers []func()
)

type Command struct {
	Version      string
	Dependencies []string
	PreRun       MambaFunc
	Login        http.HandlerFunc
	Home         http.HandlerFunc
	FAQ          http.HandlerFunc
}

func (c *Command) execute(ctx context.Context) (err error) {
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
	router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

	return nil
}

// ExecuteC executes the command.
func (c *Command) executeC(ctx context.Context) (cmd *Command, err error) {
	err = cmd.execute(ctx)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	return cmd, err
}

func (c *Command) Execute(ctx context.Context) error {
	_, err := c.executeC(ctx)
	if err != nil {
		return err
	} else {
		var srv *http.Server
		if router != nil {
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
		} else {
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
		}
		return srv.ListenAndServe()
	}
}
