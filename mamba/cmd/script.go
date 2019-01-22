package cmd

import "github.com/gofunct/mamba"

func init() {
	scriptCmd.AddCommand(dockerCmd)
}

var dockerCmd = &mamba.Command{
	Use:  "docker",
	Info: "docker scripts",
}

var scriptCmd = &mamba.Command{
	Use:  "script",
	Info: "helpful scripts",
}
