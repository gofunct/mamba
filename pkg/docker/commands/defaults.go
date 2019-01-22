package commands

import (
	"fmt"
	"github.com/gofunct/mamba"
)

func DefaultFunc() mamba.MambaFunc {
	return func(cmd *mamba.Command, args []string) {
		for _, v := range args {
			fmt.Println("arg:", v)
		}
		cmd.DebugFlags()
	}
}
