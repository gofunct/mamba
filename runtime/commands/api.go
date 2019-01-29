// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"fmt"

	"github.com/gofunct/common/pkg/logger/zap"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Select API modules which used in the service",
	Run: func(cmd *cobra.Command, args []string) {
		if V.GetBool("api.rest") || V.GetBool("api.grpc") {
			V.Set("api.enabled", true)
			if V.GetBool("api.rest") {
				V.Set("api.grpc", true)
			}
		} else {
			V.Set("api.enabled", false)
		}
		err := V.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing API configuration:", err)
		}
		fmt.Println("API configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(apiCmd)

	apiCmd.PersistentFlags().Bool("on", false, "An API modules using")
	apiCmd.PersistentFlags().Bool("rest-gateway", false, "A REST gateway module using")
	apiCmd.PersistentFlags().Bool("grpc", false, "A gRPC module using")
	zap.LogF(
		"Flag error",
		V.BindPFlag("api.enabled", apiCmd.PersistentFlags().Lookup("enabled")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("api.gateway", apiCmd.PersistentFlags().Lookup("rest-gateway")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("api.grpc", apiCmd.PersistentFlags().Lookup("grpc")),
	)
}
