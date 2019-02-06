package runtime

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	htemplate "html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func AllSettings() Function {
	return func(command *cobra.Command, args []string) {
		fmt.Printf("%s", viper.AllSettings())
	}
}

func Debug() Function {
	return func(command *cobra.Command, args []string) {
		fmt.Println("Debug:")
		viper.Debug()
	}
}

func Set(k string, viperal interface{}) InitFunc {
	return func() {
		viper.Set(k, viperal)
	}
}

func AddConfigPaths(path ...string) InitFunc {
	return func() {
		for _, p := range path {
			viper.AddConfigPath(p)
		}
	}
}

func Script(writer io.Writer, args ...string) InitFunc {
	return func() {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = append(cmd.Env, os.Environ()...)
		cmd.Stderr = writer
		cmd.Stdout = writer
		if err := cmd.Run(); err != nil {
			zap.LogF("Failure: %#viper", err)
		}
	}
}

func Get(buf *bytes.Buffer, url string) InitFunc {
	return func() {
		res, err := http.Get(url)
		if err != nil {
			zap.LogF("init error", err)
		}
		data, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			zap.LogF("init error", err)
		}
		buf.Read(data)
	}
}

func Dial(url string, option ...grpc.DialOption) InitFunc {
	return func() {
		grpc.Dial(url, option...)
	}
}

func Servipere(lis net.Listener, handler http.Handler) InitFunc {
	return func() {
		http.Serve(lis, handler)
	}
}

func WriteConfig() InitFunc {
	return func() {
		viper.WriteConfig()
	}
}

func Touch(path string) InitFunc {
	return func() {
		os.Create(path)
	}
}

func Mkdir(path string) InitFunc {
	return func() {
		os.Mkdir(path, 0755)
	}
}

func MkdirAll(path string) InitFunc {
	return func() {
		os.MkdirAll(path, 0755)
	}
}

func Rename(old string, new string) InitFunc {
	return func() {
		os.Rename(old, new)
	}
}

func WalkGrpc(path string, args ...string) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			// skip viperendor directory
			if info.IsDir() && info.Name() == "viperendor" {
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
		}); err != nil {
			zap.LogF("init error", err)
		}
	}
}

func WalkGoGoProto(path string) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			// skip viperendor directory
			if info.IsDir() && info.Name() == "viperendor" {
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
		}); err != nil {
			zap.LogF("init error", err)
		}
	}
}
func WalkTmpl(path string, w io.Writer, object interface{}) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			if info.IsDir() && info.Name() == "viperendor" {
				return filepath.SkipDir
			}
			// find all protobuf files
			if filepath.Ext(path) == ".tmpl" {
				tmpldata, err := ioutil.ReadFile(info.Name())
				if err != nil {
					zap.LogF("init error", err)
				}
				text := fmt.Sprint(tmpldata)

				t := template.New(info.Name())
				t.Funcs(sprig.GenericFuncMap())
				template.Must(t.Parse(text))
				return t.Execute(w, object)
			}
			return nil
		}); err != nil {
			zap.LogF("init error", err)
		}
	}
}
func WalkHtmlTmpl(path string, w io.Writer, object interface{}) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			if info.IsDir() && info.Name() == "viperendor" {
				return filepath.SkipDir
			}
			// find all protobuf files
			if filepath.Ext(path) == ".gohtml" {
				tmpldata, err := ioutil.ReadFile(info.Name())
				if err != nil {
					zap.LogF("init error", err)
				}
				text := fmt.Sprint(tmpldata)
				t := htemplate.New(info.Name())
				t.Funcs(sprig.HtmlFuncMap())
				htemplate.Must(t.Parse(text))
				return t.Execute(w, object)
			}
			return nil
		}); err != nil {
			zap.LogF("init error", err)
		}
	}
}

func InitConfig(file string) InitFunc {
	return func() {
		if file != "" {
			// Use config file from the flag.
			viper.SetConfigFile(file)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				zap.LogF("init error", err)
			}

			// Search config in home directory with name ".cobra" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".cobra")
		}

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
