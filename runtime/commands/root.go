package commands

import (
	"fmt"
	"github.com/gofunct/mamba/runtime/commands/app"
	"github.com/gofunct/mamba/runtime/commands/contain"
	"github.com/gofunct/mamba/runtime/commands/ctl"
	"github.com/gofunct/mamba/runtime/commands/load"
	"github.com/gofunct/mamba/runtime/commands/walk"
	"github.com/spf13/viper"

	"github.com/gofunct/common/pkg/logger/zap"

	"github.com/gofunct/mamba/runtime/config"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mamba",
	Short: "A snake powered boilerplate generator",
}

// Run adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Run() {
	zap.LogF("Bootstrap error", RootCmd.Execute())
}

func init() {
	fmt.Printf("%s %s\n\n", config.ServiceName, config.RELEASE)
	RootCmd.AddCommand(walk.RootCmd)
	RootCmd.AddCommand(load.RootCmd)
	RootCmd.AddCommand(ctl.RootCmd)
	RootCmd.AddCommand(app.RootCmd)
	RootCmd.AddCommand(contain.RootCmd)
	for _ , c := range RootCmd.Commands() {
		zap.LogE("bind pflags", viper.BindPFlags(c.Flags()))
		zap.LogE("bind pflags", viper.BindPFlags(c.PersistentFlags()))
	}
}
