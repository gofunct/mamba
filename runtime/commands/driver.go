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

var (
	databasePort   int
	databaseDriver string
)

// driverCmd represents the driver command
var driverCmd = &cobra.Command{
	Use:   "driver",
	Short: "Setup database driver settings",
	Run: func(cmd *cobra.Command, args []string) {
		err := V.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing storage driver configuration:", err)
		}
		fmt.Println("Storage driver configuration saved")
	},
}

func init() {
	storageCmd.AddCommand(driverCmd)

	if V.GetBool("storage.mysql") {
		databasePort = config.DefaultMySQLPort
		databaseDriver = config.StorageMySQL
	} else {
		databasePort = config.DefaultPostgresPort
		databaseDriver = config.StoragePostgres
	}

	driverCmd.PersistentFlags().String("host", databaseDriver, "A host name")
	driverCmd.PersistentFlags().Int("port", databasePort, "A port number")
	driverCmd.PersistentFlags().String("name", "", "A database name")
	driverCmd.PersistentFlags().StringP("username", "u", databaseDriver, "A name of database user")
	driverCmd.PersistentFlags().StringP("password", "p", databaseDriver, "An user password")
	driverCmd.PersistentFlags().Int("max-conn", 10, "Maximum available connections")
	driverCmd.PersistentFlags().Int("idle-conn", 1, "Count of idle connections")
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.config.host", driverCmd.PersistentFlags().Lookup("host")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.config.port", driverCmd.PersistentFlags().Lookup("port")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.config.name", driverCmd.PersistentFlags().Lookup("name")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.config.username", driverCmd.PersistentFlags().Lookup("username")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.config.password", driverCmd.PersistentFlags().Lookup("password")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.config.connections.max", driverCmd.PersistentFlags().Lookup("max-conn")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("storage.config.connections.idle", driverCmd.PersistentFlags().Lookup("idle-conn")),
	)
}
