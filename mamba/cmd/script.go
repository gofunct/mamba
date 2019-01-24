package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	scriptCmd.AddCommand(dockerCmd)
}

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "docker scripts",
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "helpful scripts",
}
