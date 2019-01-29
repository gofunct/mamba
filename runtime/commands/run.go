package commands

import (
	"fmt"

	"github.com/gofunct/common/pkg/logger/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the api command
var runCmd = &cobra.Command{
	Use: "run",

	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("api.rest") || viper.GetBool("api.grpc") {
			viper.Set("api.enabled", true)
			if viper.GetBool("api.rest") {
				viper.Set("api.grpc", true)
			}
		} else {
			viper.Set("api.enabled", false)
		}
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing API configuration:", err)
		}
		fmt.Println("API configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(apiCmd)

	apiCmd.PersistentFlags().Bool("enabled", false, "An API modules using")
	apiCmd.PersistentFlags().Bool("rest-gateway", false, "A REST gateway module using")
	apiCmd.PersistentFlags().Bool("grpc", false, "A gRPC module using")
	zap.LogF(
		"Flag error",
		viper.BindPFlag("api.enabled", apiCmd.PersistentFlags().Lookup("enabled")),
	)
	zap.LogF(
		"Flag error",
		viper.BindPFlag("api.gateway", apiCmd.PersistentFlags().Lookup("rest-gateway")),
	)
	zap.LogF(
		"Flag error",
		viper.BindPFlag("api.grpc", apiCmd.PersistentFlags().Lookup("grpc")),
	)
}
