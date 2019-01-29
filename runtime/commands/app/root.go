package app

import (
	"fmt"
	"github.com/gofunct/mamba/runtime/input"
	"os"

	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/gofunct/mamba/runtime/generator"

	"github.com/gofunct/mamba/runtime/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	V       = config.Viper
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "🐍 Generate boilerplate to bootstrap an application",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Configuration
		if err := V.Unmarshal(&cfg); err != nil {
			fmt.Println("Error parsing of configuration, used default:", err)
		}
		cfg = input.Inquire(cfg)
		if cfg.Storage.MySQL &&
			cfg.Storage.Config.Port == config.DefaultPostgresPort {
			cfg.Storage.Config.Driver = config.StorageMySQL
			cfg.Storage.Config.Host = config.StorageMySQL
			cfg.Storage.Config.Port = config.DefaultMySQLPort
			cfg.Storage.Config.Username = config.StorageMySQL
			cfg.Storage.Config.Password = config.StorageMySQL
		}
		if cfg.Storage.Postgres &&
			cfg.Storage.Config.Port == config.DefaultMySQLPort {
			cfg.Storage.Config.Driver = config.StoragePostgres
			cfg.Storage.Config.Host = config.StoragePostgres
			cfg.Storage.Config.Port = config.DefaultPostgresPort
			cfg.Storage.Config.Username = config.StoragePostgres
			cfg.Storage.Config.Password = config.StoragePostgres
		}
		generator.Run(cfg)
	},
}

func init() {
	fmt.Printf("%s %s\n\n", config.ServiceName, config.RELEASE)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", os.Getenv("HOME")+"/.mamba.yaml", "config file (default is ~/.mamba.yaml)")
	RootCmd.PersistentFlags().String("templates", "templates", "templates dir")
	RootCmd.PersistentFlags().String("service", "", "A boilerplate service repository dir")

	RootCmd.AddCommand(
		apiCmd,
		debugCmd,
		gkeCmd,
		storageCmd,
		newCmd,
		)
	zap.LogF(
		"Flag error",
		V.BindPFlag("directories.templates", RootCmd.PersistentFlags().Lookup("templates")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("directories.service", RootCmd.PersistentFlags().Lookup("service")),
	)
}
