package commands

import (
	"fmt"
	"github.com/gofunct/mamba/function"
	"github.com/spf13/cobra"
)

func DefaultFunc() function.RunFunc {
	return func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			fmt.Println("arg:", v)
		}
		cmd.DebugFlags()
	}
}
