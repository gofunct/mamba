// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package app

import (
	"fmt"
	"github.com/gofunct/common/pkg/logger/zap"

	"github.com/spf13/cobra"
)

// configCmd represents API settings command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "🐍  Setup API settings",
	Run: func(cmd *cobra.Command, args []string) {
		err := V.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing API settings:", err)
		}
		fmt.Println("API configuration saved")
	},
}

func init() {
	configCmd.PersistentFlags().Int("port", 8000, "A service port number")
	configCmd.PersistentFlags().Int("gateway-port", 8480, "A service rest gateway port number")
	zap.LogF("Flag error", V.BindPFlag("api.config.port", configCmd.PersistentFlags().Lookup("port")))
	zap.LogF(
		"Flag error",
		V.BindPFlag("api.config.gateway.port", configCmd.PersistentFlags().Lookup("gateway-port")),
	)
}
