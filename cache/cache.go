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

func init() {
	Cache.v = viper.GetViper()
}

var (
	OSFS  = afero.NewOsFs()
	Cache *cache
)

type cache struct {
	v *viper.Viper
}

func (c *cache) Bind(cmd *cobra.Command) error {
	{
		c.v.SetFs(OSFS)
		c.v.SetConfigName("." + filepath.Base(os.Getenv("PWD")))
		c.v.AddConfigPath(os.Getenv("HOME"))
		c.v.AutomaticEnv()
		c.v.AllowEmptyEnv(true)
		c.v.SetTypeByDefaultValue(true)
		c.v.SetDefault("env.base", filepath.Base(os.Getenv("PWD")))
		home, _ := os.LookupEnv("HOME")
		c.v.SetDefault("env.home", home)
		gopath, _ := os.LookupEnv("GOPATH")
		c.v.SetDefault("env.gopath", gopath)
		user, _ := os.LookupEnv("USER")
		c.v.SetDefault("env.user", user)
		modules, _ := os.LookupEnv("GO111MODULES")
		c.v.SetDefault("env.modules", modules)
		creds, _ := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
		c.v.SetDefault("env.creds", creds)
		pwd, _ := os.LookupEnv("PWD")
		c.v.SetDefault("env.absolute", pwd)
		host, _ := os.Hostname()
		c.v.SetDefault("env.host", host)
	}

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
	if err := c.Write(); err != nil {
		return err
	}
	cmd.AddCommand(c.DebugCmd())
	return nil
}

func (c *cache) Write() error {
	// If a config file is found, read it in.
	b, err := afero.Exists(OSFS, os.Getenv("HOME")+"/."+filepath.Base(os.Getenv("PWD")+".json"))
	if err != nil {
		return errors.WithStack(err)
	}
	if !b {
		f, err := os.Create(os.Getenv("HOME") + "/." + filepath.Base(os.Getenv("PWD")+".json"))
		if err != nil {
			return errors.WithStack(err)
		}
		c.v.SetConfigFile(f.Name())
	}
	if err := c.v.ReadInConfig(); err != nil {
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

func (c *cache) DebugCmd() *cobra.Command {
	// debugCmd represents the debug command
	var debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Debug your current configuration settings",
		Run: func(cmd *cobra.Command, args []string) {
			c.v.Debug()
		},
	}
	return debugCmd
}
