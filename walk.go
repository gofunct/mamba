package mamba

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/mamba/pkg/function"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func (c *Command) WalkGrpc(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			args = []string{
				"-I=.",
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src")),
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gogo", "protobuf", "protobuf")),
				fmt.Sprintf("--proto_path=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com")),
				"--go_out=plugins=grpc:.",
				path,
			}
			cmd := exec.Command("protoc", args...)
			err = cmd.Run()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (c *Command) WalkGoGoProto(path string) error {
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
			err = cmd.Run()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *Command) WalkTmpl(wr io.Writer, text string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		t := template.New("top")
		t.Funcs(sprig.GenericFuncMap())
		template.Must(t.Parse(text))
		return t.Execute(wr, viper.AllSettings())
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".tmpl" {

		}
		return nil
	}
}

func (c *Command) WalkShell(args ...string) filepath.WalkFunc {
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

func (c *Command) WalkGo(args ...string) filepath.WalkFunc {
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

func (c *Command) WalkMakefile(args ...string) filepath.WalkFunc {
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

func (c *Command) WalkDockerfile(args ...string) filepath.WalkFunc {
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
