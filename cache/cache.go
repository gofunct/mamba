package cache

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func init() {
	fmt.Println("setting up cache...")
	{
		Cache.fmap = sprig.GenericFuncMap()
	}
	{
		Cache.v = viper.GetViper()
		Cache.v.SetFs(osFs)
		Cache.v.SetConfigName("." + filepath.Base(os.Getenv("PWD")))
		Cache.v.AddConfigPath(os.Getenv("HOME"))
		Cache.v.AutomaticEnv()
		Cache.v.AllowEmptyEnv(true)
		Cache.v.SetTypeByDefaultValue(true)
	}

	{
		Cache.v.SetDefault("env.base", filepath.Base(os.Getenv("PWD")))
		home, _ := os.LookupEnv("HOME")
		Cache.v.SetDefault("env.home", home)
		gopath, _ := os.LookupEnv("GOPATH")
		Cache.v.SetDefault("env.gopath", gopath)
		user, _ := os.LookupEnv("USER")
		Cache.v.SetDefault("env.user", user)
		modules, _ := os.LookupEnv("GO111MODULES")
		Cache.v.SetDefault("env.modules", modules)
		creds, _ := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
		Cache.v.SetDefault("env.creds", creds)
		pwd, _ := os.LookupEnv("PWD")
		Cache.v.SetDefault("env.absolute", pwd)
		host, _ := os.Hostname()
		Cache.v.SetDefault("env.host", host)
		Cache.v.SetDefault("env.base", filepath.Base(os.Getenv("PWD")))
		home, _ = os.LookupEnv("HOME")
		Cache.v.SetDefault("env.home", home)
		gopath, _ = os.LookupEnv("GOPATH")
		Cache.v.SetDefault("env.gopath", gopath)
		user, _ = os.LookupEnv("USER")
		Cache.v.SetDefault("env.user", user)
		modules, _ = os.LookupEnv("GO111MODULES")
		Cache.v.SetDefault("env.modules", modules)
		creds, _ = os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
		Cache.v.SetDefault("env.creds", creds)
		pwd, _ = os.LookupEnv("PWD")
		Cache.v.SetDefault("env.absolute", pwd)
		host, _ = os.Hostname()
		Cache.v.SetDefault("env.host", host)
	}
	if err := Cache.Write(); err != nil {
		fmt.Printf("failed to write config from cache %s\n", err.Error())
	}
}

var (
	osFs  = afero.NewOsFs()
	Cache = new(cache)
)

type cache struct {
	v    *viper.Viper
	fmap template.FuncMap
}

func (c *cache) Write() error {
	// If a config file is found, read it in.
	b, err := afero.Exists(osFs, os.Getenv("HOME")+"/."+filepath.Base(os.Getenv("PWD")+".json"))
	if err != nil {
		return errors.WithStack(err)
	}
	if !b {
		f, err := os.Create(os.Getenv("HOME") + "/." + filepath.Base(os.Getenv("PWD")+".json"))
		if err != nil {
			return errors.WithStack(err)
		}
		Cache.v.SetConfigFile(f.Name())
	}
	if err := Cache.v.ReadInConfig(); err != nil {
		log.Println("failed to read config file, writing defaults...")
		if err := c.v.WriteConfig(); err != nil {
			return errors.Wrap(err, "failed to write config")
		}
		log.Println("Using config file:", c.v.ConfigFileUsed())
		if err := c.v.WriteConfig(); err != nil {
			return errors.WithStack(err)
		}

	} else {
		log.Println("Using config file:", c.v.ConfigFileUsed())
		if err := c.v.WriteConfig(); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *cache) WrapCobra(cmd *cobra.Command) error {
	// debugCmd represents the debug command
	var debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Debug your current configuration settings",
		Run: func(cmd *cobra.Command, args []string) {
			c.v.Debug()
		},
	}
	cmd.AddCommand(debugCmd)

	if err := c.v.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := c.v.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}
	c.v.SetDefault(cmd.Name()+".meta", cmd.Annotations)

	for _, cmds := range cmd.Commands() {
		if err := c.v.BindPFlags(cmds.Flags()); err != nil {
			return err
		}
		if err := c.v.BindPFlags(cmds.PersistentFlags()); err != nil {
			return err
		}
		c.v.SetDefault(cmds.Name()+".meta", cmd.Annotations)
	}
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := c.Write(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (c *cache) GetFs() afero.Fs {
	return osFs
}
