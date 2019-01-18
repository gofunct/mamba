package function

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

func GrpcWalkFunc(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			args = append([]string{"protoc", "--go_out=plugins=grpc:."}, args...)
			ValidateString(args...)
		}
		return nil
	}
}

func GoGoWalkFunc(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".proto" {
			args = append([]string{
				"protoc",
				"-I=.",
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src")),
				fmt.Sprintf("-I=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gogo", "protobuf", "protobuf")),
				fmt.Sprintf("--proto_path=%s", filepath.Join(os.Getenv("GOPATH"), "src", "github.com")),
				"--gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:.",
			}, args...)
			ValidateString(args...)
		}
		return nil
	}
}

func TmplWalkFunc() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".tmpl" {

		}
		return nil
	}
}

func ShellWalkFunc(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if filepath.Ext(path) == ".sh" {
			args = append([]string{"bash"}, args...)
			ValidateString(args...)
		}
		return nil
	}
}

func ConfigWalkFunc() filepath.WalkFunc {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetConfigName(".mamba")
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		if info.Name() == ".mamba" {
			viper.AddConfigPath(path)
		}

		return nil
	}
}

func GoWalkFunc(args ...string) filepath.WalkFunc {
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
			ValidateString(args...)
		}
		return nil
	}
}

func MakefileWalkFunc(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if !info.IsDir() && info.Name() == "Makefile" {
			args = append([]string{"make"}, args...)
			ValidateString(args...)
		}
		return nil
	}
}

func DockerfileWalkFunc(args ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		// find all protobuf files
		if !info.IsDir() && info.Name() == "Dockerfile" {
			args = append([]string{"go"}, args...)
			ValidateString(args...)
		}

		return nil
	}
}
