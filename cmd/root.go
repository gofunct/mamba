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
	kitlog "github.com/go-kit/kit/log"
	"github.com/gofunct/mamba/cmd/gcloud"
	"github.com/gofunct/mamba/cmd/local"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var (
	fs = afero.NewOsFs()
	in 	string
	out		string
)


func init() {
	{
		rootCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "path to input directory")
		rootCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "path to output directory")
	}
	{
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
		log.SetOutput(kitlog.NewStdlibAdapter(logger))
	}

	{
		rootCmd.AddCommand(gcloud.RootCmd)
		rootCmd.AddCommand(local.RootCmd)
		rootCmd.AddCommand(protocCmd)
		rootCmd.AddCommand(htmlCmd)
		rootCmd.AddCommand(testCmd)
		rootCmd.AddCommand(protocGenCmd)
		rootCmd.AddCommand(debugCmd)
	}
	{
		if err := Bind(rootCmd, local.RootCmd, protocGenCmd, protocCmd, htmlCmd, testCmd, debugCmd); err != nil {
			log.Println("failed to bind config to commands\n", err.Error())
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "mamba",
	Short: "A general purpose scripting utility for developers and administrators",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Bind(c ...*cobra.Command) error {
	{
		viper.SetFs(fs)
		viper.SetConfigName("."+filepath.Base(os.Getenv("PWD")))
		viper.AddConfigPath(os.Getenv("HOME"))
		viper.AutomaticEnv()
		viper.AllowEmptyEnv(true)
		viper.SetTypeByDefaultValue(true)
		viper.SetDefault("env.base", filepath.Base(os.Getenv("PWD")))
		home, _ := os.LookupEnv("HOME")
		viper.SetDefault("env.home", home)
		gopath, _ := os.LookupEnv("GOPATH")
		viper.SetDefault("env.gopath", gopath)
		user, _ := os.LookupEnv("USER")
		viper.SetDefault("env.user", user)
		modules, _ := os.LookupEnv("GO111MODULES")
		viper.SetDefault("env.modules", modules)
		creds, _ := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
		viper.SetDefault("env.creds", creds)
		pwd, _ := os.LookupEnv("PWD")
		viper.SetDefault("env.absolute", pwd)
		host, _ := os.Hostname()
		viper.SetDefault("env.host", host)
	}


	for _, cmd := range c {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		viper.SetDefault(cmd.Name()+".meta", cmd.Annotations)
	}
	if err := write(); err != nil {
		return err
	}
	return nil
}


func write() error {
	// If a config file is found, read it in.
	b, err := afero.Exists(fs, os.Getenv("HOME")+"/."+filepath.Base(os.Getenv("PWD")+".json"))
	if err != nil {
		return errors.WithStack(err)
	}
	if !b {
		f, err := os.Create(os.Getenv("HOME")+"/."+filepath.Base(os.Getenv("PWD")+".json"))
		if err != nil {
			return errors.WithStack(err)
		}
		viper.SetConfigFile(f.Name())
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Println("failed to read config file, writing defaults...")
		if err := viper.WriteConfig(); err != nil {
			return errors.Wrap(err, "failed to write config")
		}
		log.Println("Using config file:", viper.ConfigFileUsed())
		if err := viper.WriteConfig(); err != nil {
			return errors.WithStack(err)
		}

	} else {
		log.Println("Using config file:", viper.ConfigFileUsed())
		if err := viper.WriteConfig(); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}