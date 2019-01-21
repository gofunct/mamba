package docker

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(CreateCmd)
	RootCmd.AddCommand(imageCmd())

}

var RootCmd = &cobra.Command{
	Use:   "docker",
	Short: "docker scripts",
}
