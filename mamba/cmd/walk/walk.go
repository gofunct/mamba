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
	"github.com/gofunct/mamba/runtime/function"
	"github.com/pkg/errors"
	"github.com/shiyanhui/hero"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	walkCmd.AddCommand(goGoCmd, grpcCmd, htmlCmd, jsCmd)
}

// protocCmd represents the protoc command
var goGoCmd = &cobra.Command{
	Use:   "gogo",
	Short: "🐍 Compile gogo protobuf stubs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := WalkGoGoProto(in); err != nil {
			logger.Fatalf("%s", errors.WithStack(err))
		}
	},
}

// protocCmd represents the protoc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "🐍 Compile grpc protobuf stubs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := WalkGrpc(in); err != nil {
			logger.Fatalf("%s", errors.WithStack(err))
		}
	},
}

// protocCmd represents the protoc command
var jsCmd = &cobra.Command{
	Use:   "js",
	Short: "🐍 Compile grpc javascript protobuf stubs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := WalkProtoJs(in); err != nil {
			logger.Fatalf("%s", errors.WithStack(err))
		}
	},
}

// protocCmd represents the protoc command
var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "🐍 Walk a filepath with a given function and file extension",
}

var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "🐍 Generate html files",
	Run: func(cmd *cobra.Command, args []string) {
		generate(in, out, pkg)
	},
}

func generate(source, dest, pkg string) {
	hero.Generate(source, dest, pkg)
}

func WalkGrpc(path string) error {

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			// args
			args := []string{
				"-I=.",
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("PWD"), "third_party")),
				fmt.Sprintf("--go_out=plugins=grpc:%s", "."),
				path,
			}

			cmd := exec.Command("protoc", args...)
			cmd.Env = os.Environ()
			o, _ := cmd.Output()
			fmt.Println(o)
		}
		return nil
	})
}

func WalkProtoJs(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			args := []string{
				"-I=.",
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src")),
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("PWD"), "third_party")),
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gogo", "protobuf", "protobuf")),
				fmt.Sprintf("--proto_path=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com")),
				"--gopherjs_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/johanbrandhorst/protobuf/ptypes/timestamp:.",
				path,
			}
			cmd := exec.Command("protoc", args...)
			cmd.Env = os.Environ()
			fmt.Println("cmd.path:", cmd.Path, "cmd.dir:", cmd.Dir, "cmd.args", cmd.Args, "cmd.state", cmd.ProcessState.String())
			o, _ := cmd.Output()
			fmt.Println(o)
		}
		return nil
	})
}

func WalkGoGoProto(path string) error {

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			// args
			args := []string{
				"-I=.",
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src")),
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gogo", "protobuf", "protobuf")),
				fmt.Sprintf("--proto_path=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com")),
				"--gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:.",
				path,
			}
			cmd := exec.Command("protoc", args...)
			cmd.Env = os.Environ()
			o, _ := cmd.Output()
			fmt.Println(o)
		}
		return nil
	})
}
func WalkTmpl(wr io.Writer, text string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".tmpl" {

		}
		return nil
	}
}

func WalkShell(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".sh" {
			args = append([]string{"bash"}, args...)
			function.ValidateString(args...)
		}
		return nil
	}
}

func WalkGo(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if path == "" || info == nil || err != nil {
			log.Fatalf("Walkfunc failure: %s %v %s", path, info, err)
		}
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".go" {
			args = append([]string{"go"}, args...)
			function.ValidateString(args...)
		}
		return nil
	}
}

func WalkMakefile(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if !info.IsDir() && info.Name() == "Makefile" {
			args = append([]string{"make"}, args...)
			function.ValidateString(args...)
		}
		return nil
	}
}

func WalkDockerfile(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if !info.IsDir() && info.Name() == "Dockerfile" {
			args = append([]string{"docker"}, args...)
			function.ValidateString(args...)
		}

		return nil
	}
}
