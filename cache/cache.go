package cache

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var (
	OSFS = afero.NewOsFs()
)

func Bind(c ...*cobra.Command) error {
	{
		viper.SetFs(OSFS)
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
	b, err := afero.Exists(OSFS, os.Getenv("HOME")+"/."+filepath.Base(os.Getenv("PWD")+".json"))
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