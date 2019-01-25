// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
)

// getCmd represents the get command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "download files from remote or local sources",
	Run: func(cmd *cobra.Command, args []string) {
		var mode = getter.ClientModeAny

		// Get the pwd
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting wd: %s", err)
		}

		opts := []getter.ClientOption{}
		opts = append(opts, getter.WithProgress(defaultProgressBar))

		ctx, cancel := context.WithCancel(context.Background())
		// Build the client
		client := &getter.Client{
			Ctx:     ctx,
			Src:     args[0],
			Dst:     args[1],
			Pwd:     pwd,
			Mode:    mode,
			Options: opts,
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		errChan := make(chan error, 2)
		go func() {
			defer wg.Done()
			defer cancel()
			if err := client.Get(); err != nil {
				errChan <- err
			}
		}()

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)

		select {
		case sig := <-c:
			signal.Reset(os.Interrupt)
			cancel()
			wg.Wait()
			log.Printf("signal %v", sig)
		case <-ctx.Done():
			wg.Wait()
			log.Printf("success!")
		case err := <-errChan:
			wg.Wait()
			log.Fatalf("Error downloading: %s", err)
		}
	},
}

var defaultProgressBar getter.ProgressTracker = &ProgressBar{}

type ProgressBar struct {
	// lock everything below
	lock sync.Mutex

	pool *pb.Pool

	pbs int
}

func ProgressBarConfig(bar *pb.ProgressBar, prefix string) {
	bar.SetUnits(pb.U_BYTES)
	bar.Prefix(prefix)
}

func (cpb *ProgressBar) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	newPb := pb.New64(totalSize)
	newPb.Set64(currentSize)
	ProgressBarConfig(newPb, filepath.Base(src))
	if cpb.pool == nil {
		cpb.pool = pb.NewPool()
		cpb.pool.Start()
	}
	cpb.pool.Add(newPb)
	reader := newPb.NewProxyReader(stream)

	cpb.pbs++
	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			newPb.Finish()
			cpb.pbs--
			if cpb.pbs <= 0 {
				cpb.pool.Stop()
				cpb.pool = nil
			}
			return nil
		},
	}
}

type readCloser struct {
	io.Reader
	close func() error
}

func (c *readCloser) Close() error { return c.close() }
