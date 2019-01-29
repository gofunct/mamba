// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package app

import (
	"fmt"

	"github.com/gofunct/common/pkg/logger/zap"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "🐍  Select API modules which used in the service",
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
	apiCmd.AddCommand(configCmd)
	apiCmd.Flags().Bool("enabled", false, "An API modules using")
	apiCmd.Flags().Bool("rest-gateway", false, "A REST gateway module using")
	apiCmd.Flags().Bool("grpc", false, "A gRPC module using")
	zap.LogF(
		"Flag error",
		V.BindPFlag("api.enabled", apiCmd.Flags().Lookup("enabled")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("api.gateway", apiCmd.Flags().Lookup("rest-gateway")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("api.grpc", apiCmd.Flags().Lookup("grpc")),
	)
}
