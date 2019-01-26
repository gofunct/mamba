package function

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gofunct/mamba/runtime/input"
	"github.com/mitchellh/go-homedir"
	"github.com/robfig/cron"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	htemplate "html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"text/template"

	kitlog "github.com/go-kit/kit/log"
	"os"
	"os/exec"
)

func init() {
	fs = afero.NewOsFs()
	v = viper.GetViper()
	q = input.DefaultUI()
	v.SetFs(fs)
	v.SetConfigName("mamba.json")
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)
}

type CobraFunc func(command *cobra.Command, args []string)
type InitFunc func()

var v *viper.Viper
var fs afero.Fs
var q *input.UI

type Scripter struct {
	ProjectId    string
	Initializers []func()
	Run          CobraFunc
	PostRun      CobraFunc
	cron         *cron.Cron
}

func (s *Scripter) Execute(w io.Writer, id, name, info string) error {
	cobra.OnInitialize(s.Initializers...)

	if s.cron != nil {
		go func() { s.cron.Start() }()
	}

	cmd := &cobra.Command{
		Use:   name,
		Short: info,
	}
	cmd.Run = s.Run
	cmd.PostRun = s.PostRun
	cmd.SetOutput(w)
	if !cmd.Runnable() {
		return errors.New("command is not runnable")
	}
	return cmd.Execute()
}

func (s *Scripter) AddJob(sched string, f InitFunc) {
	if s.cron == nil {
		s.cron = cron.New()
	}
	if err := s.cron.AddFunc(sched, f); err != nil {
		ERR(err)
	}
}

func AllSettings() CobraFunc {
	return func(command *cobra.Command, args []string) {
		fmt.Printf("%s", v.AllSettings())
	}
}

func Debug() CobraFunc {
	return func(command *cobra.Command, args []string) {
		fmt.Println("Debug:")
		v.Debug()
	}
}

func Set(k string, val interface{}) InitFunc {
	return func() {
		v.Set(k, val)
	}
}

func (s *Scripter) Get(k string) interface{} {
	return v.Get(k)
}

func (s *Scripter) Unmarshal(object interface{}) InitFunc {
	return func() {
		if err := v.Unmarshal(object); err != nil {
			log.Fatalf("Failure: %#v", err)
		}
	}
}

func AddConfigPaths(path ...string) InitFunc {
	return func() {
		for _, p := range path {
			v.AddConfigPath(p)
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
			log.Fatalf("Failure: %#v", err)
		}
	}
}

func Get(buf *bytes.Buffer, url string) InitFunc {
	return func() {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		buf.Read(data)
	}
}

func Dial(url string, option ...grpc.DialOption) InitFunc {
	return func() {
		grpc.Dial(url, option...)
	}
}

func Serve(lis net.Listener, handler http.Handler) InitFunc {
	return func() {
		http.Serve(lis, handler)
	}
}

func WriteConfig() InitFunc {
	return func() {
		v.WriteConfig()
	}
}

func Touch(path string) InitFunc {
	return func() {
		fs.Create(path)
	}
}

func Mkdir(path string) InitFunc {
	return func() {
		fs.Mkdir(path, 0755)
	}
}

func MkdirAll(path string) InitFunc {
	return func() {
		fs.MkdirAll(path, 0755)
	}
}

func RemoveAll(path string) InitFunc {
	return func() {
		fs.RemoveAll(path)
	}
}

func Rename(old string, new string) InitFunc {
	return func() {
		fs.Rename(old, new)
	}
}

func WalkGrpc(path string, args ...string) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
		}); err != nil {
			ERR(err)
		}
	}
}

func WalkGoGoProto(path string) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
		}); err != nil {
			ERR(err)
		}
	}
}
func WalkTmpl(path string, w io.Writer, object interface{}) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			if info.IsDir() && info.Name() == "vendor" {
				return filepath.SkipDir
			}
			// find all protobuf files
			if filepath.Ext(path) == ".tmpl" {
				tmpldata, err := ioutil.ReadFile(info.Name())
				if err != nil {
					ERR(err)
				}
				text := fmt.Sprint(tmpldata)

				t := template.New(info.Name())
				t.Funcs(sprig.GenericFuncMap())
				template.Must(t.Parse(text))
				return t.Execute(w, object)
			}
			return nil
		}); err != nil {
			ERR(err)
		}
	}
}
func WalkHtmlTmpl(path string, w io.Writer, object interface{}) InitFunc {
	return func() {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			if info.IsDir() && info.Name() == "vendor" {
				return filepath.SkipDir
			}
			// find all protobuf files
			if filepath.Ext(path) == ".gohtml" {
				tmpldata, err := ioutil.ReadFile(info.Name())
				if err != nil {
					ERR(err)
				}
				text := fmt.Sprint(tmpldata)
				t := htemplate.New(info.Name())
				t.Funcs(sprig.HtmlFuncMap())
				htemplate.Must(t.Parse(text))
				return t.Execute(w, object)
			}
			return nil
		}); err != nil {
			ERR(err)
		}
	}
}

func ReplaceLogger(writer io.Writer) InitFunc {
	return func() {
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(writer))
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)
		log.SetOutput(kitlog.NewStdlibAdapter(logger))
	}
}

func ERR(err error) {
	log.Fatalf("%#v\n%s\n", err, err.Error())
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
				ERR(err)
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
