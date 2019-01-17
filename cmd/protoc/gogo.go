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

package protoc

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	dir string
)

// protocCmd represents the protoc command
var GoGoCmd = &cobra.Command{
	Use:   "gogo",
	Short: "generate protocol buffers",
	Run: func(cmd *cobra.Command, args []string) {
		if err := GoGo(dir); err != nil {
			log.Fatalln("failed to execute command", errors.WithStack(err))
		}
	},
}

func init() {
	ProtocCmd.AddCommand(GoGoCmd)
	GoGoCmd.Flags().StringVar(&dir, "dir", "", "path to directory containing protobuf files")
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC, "exec", kitlog.DefaultCaller, "dir", dir)
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
}

func GoGo(d string) error {

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			// args
			args := []string{
				"-I=.",
				"-I=vendor/github.com/grpc-ecosystem/grpc-gateway",
				"-I=vendor/github.com/gogo/googleapis",
				"-I=vendor",
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src")),
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gogo", "protobuf", "protobuf")),
				fmt.Sprintf("--proto_path=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com")),
				"--gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types," +
					"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types," +
					"Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types," +
					"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types," +
					"Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types," +
					"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types" +
					"Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api" +
					"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types" +
					"plugins=grpc+embedded:.",
				"gogo_" + path,
			}
			cmd := exec.Command("protoc", args...)
			log.Print("starting command")
			cmd.Env = os.Environ()
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	})
}
