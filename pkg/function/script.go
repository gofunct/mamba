package function

import (
	"errors"
	"fmt"
	"github.com/gofunct/mamba/pkg/input"
	"github.com/prometheus/common/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"

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
	v.SetEnvPrefix("mamba")

}
type CobraFunc func(command *cobra.Command, args []string)
type InitFunc func()

var v *viper.Viper
var fs afero.Fs
var q 	*input.UI

type Scripter struct {
	ProjectId string
	Initializers []func()
	PreRun CobraFunc
	Run CobraFunc
	PostRun CobraFunc
	Output io.Writer
}

func (s *Scripter) Execute(id, name, info string) error {
	cobra.OnInitialize(s.Initializers...)
	cmd := &cobra.Command{
		Use: name,
		Short: info,
	}
	cmd.PreRun = s.PreRun
	cmd.Run = s.Run
	cmd.PreRun = s.PostRun
	cmd.SetOutput(s.Output)
	if !cmd.Runnable() {
		return errors.New("command is not runnable")
	}
	return cmd.Execute()
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

func (s *Scripter) Get(k string)  interface{} {
		return v.Get(k)
}

func (s *Scripter) Unmarshal(object interface{}) InitFunc {
	return func() {
		if err := v.Unmarshal(object); err != nil {
			log.Fatalf("Failure: %#v", err)
		}
	}
}

func (s *Scripter) AddConfigPaths(path ...string) InitFunc {
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
