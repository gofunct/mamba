// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"fmt"

	"github.com/gofunct/common/pkg/logger/zap"

	"github.com/spf13/cobra"
)

// gkeCmd represents the GKE command
var gkeCmd = &cobra.Command{
	Use:   "GKE",
	Short: "Setup Google Kubernetes Engine properties to deploy the service",
	Run: func(cmd *cobra.Command, args []string) {
		err := V.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing GKE configuration:", err)
		}
		fmt.Println("GKE configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(gkeCmd)

	gkeCmd.PersistentFlags().Bool("enabled", false, "A Google Kubernetes Engine enabled")
	gkeCmd.PersistentFlags().String("project", "my-project-id", "A project ID in GCP")
	gkeCmd.PersistentFlags().String("zone", "europe-west1-b", "A compute zone in GCP")
	gkeCmd.PersistentFlags().String("cluster", "my-cluster-name", "A cluster name in GKE")
	zap.LogF(
		"Flag error",
		V.BindPFlag("gke.enabled", gkeCmd.PersistentFlags().Lookup("enabled")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("gke.project", gkeCmd.PersistentFlags().Lookup("project")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("gke.zone", gkeCmd.PersistentFlags().Lookup("zone")),
	)
	zap.LogF(
		"Flag error",
		V.BindPFlag("gke.cluster", gkeCmd.PersistentFlags().Lookup("cluster")),
	)
}
