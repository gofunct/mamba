// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"fmt"
	"github.com/gofunct/common/pkg/logger/zap"

	"github.com/gofunct/mamba/runtime/config"
	"github.com/spf13/cobra"
)

// storageCmd represents the storage command
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Setup your storage modules",
	Run: func(cmd *cobra.Command, args []string) {
		if V.GetBool("storage.postgres") || V.GetBool("storage.mysql") {
			V.Set("storage.enabled", true)
			if V.GetBool("storage.postgres") {
				V.Set("storage.mysql", false)
			}
			if V.GetBool("storage.mysql") &&
				V.GetInt("storage.driver.port") == config.DefaultPostgresPort {
				V.Set("storage.driver.host", config.StorageMySQL)
				V.Set("storage.driver.port", config.DefaultMySQLPort)
				V.Set("storage.driver.name", config.StorageMySQL)
				V.Set("storage.driver.username", config.StorageMySQL)
				V.Set("storage.driver.password", config.StorageMySQL)
			}
			if V.GetBool("storage.postgres") &&
				V.GetInt("storage.driver.port") == config.DefaultMySQLPort {
				V.Set("storage.driver.host", config.StoragePostgres)
				V.Set("storage.driver.port", config.DefaultPostgresPort)
				V.Set("storage.driver.name", config.StoragePostgres)
				V.Set("storage.driver.username", config.StoragePostgres)
				V.Set("storage.driver.password", config.StoragePostgres)
			}
		} else {
			V.Set("storage.enabled", false)
		}
		err := V.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing storage configuration:", err)
		}
		fmt.Println("Storage configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(storageCmd)

	storageCmd.PersistentFlags().Bool("enabled", false, "A Storage modules using")
	storageCmd.PersistentFlags().Bool("postgres", false, "A postgres module using")
	storageCmd.PersistentFlags().Bool("mysql", false, "A mysql module using")
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.enabled", storageCmd.PersistentFlags().Lookup("enabled")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.postgres", storageCmd.PersistentFlags().Lookup("postgres")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.mysql", storageCmd.PersistentFlags().Lookup("mysql")),
	)
}
