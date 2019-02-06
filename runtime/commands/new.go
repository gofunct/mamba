// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/gofunct/mamba/runtime/config"
	"github.com/gofunct/mamba/runtime/generator"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generates new service from templates using default/config settings",
	Long: `In this mode, you'll be not asked about everything.
The configuration file will be used for all other data,
such as the host, port, etc., if you have saved it before.
Otherwise, the default settings will be used.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := new(config.Config)
		if err := V.Unmarshal(&cfg); err != nil {
			fmt.Println("Error parsing of configuration, used default:", err)
		}
		if !path.IsAbs(cfg.Directories.Templates) {
			if currentDir, err := os.Getwd(); err == nil {
				cfg.Directories.Templates = path.Join(currentDir, cfg.Directories.Templates)
			}
		}
		if cfg.Directories.Service == "" {
			if currentDir, err := os.Getwd(); err == nil {
				cfg.Directories.Service = path.Join(path.Dir(currentDir), cfg.Project.Name)
			}
		}
		generator.Run(cfg)
	},
}

func init() {
	RootCmd.AddCommand(newCmd)

	newCmd.PersistentFlags().String("name", "my-service", "A name of your new service")
	newCmd.PersistentFlags().String("description", "My service", "A description of your new service")
	RootCmd.PersistentFlags().String("github", "my-account", "A Github account name")
	RootCmd.PersistentFlags().Bool("git-init", false, "Initialize repository with git")
	RootCmd.PersistentFlags().Bool("contract-example", false, "A example of contract API using")
	zap.LogF("Flag error", V.BindPFlag("name", newCmd.PersistentFlags().Lookup("name")))
	zap.LogF("Flag error", V.BindPFlag("description", newCmd.PersistentFlags().Lookup("description")))
	zap.LogF("Flag error", V.BindPFlag("github", RootCmd.PersistentFlags().Lookup("github")))
	zap.LogF("Flag error", V.BindPFlag("gitinit", RootCmd.PersistentFlags().Lookup("git-init")))
	zap.LogF("Flag error", V.BindPFlag("contract", RootCmd.PersistentFlags().Lookup("contract-example")))
}
