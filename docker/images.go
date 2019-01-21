package docker

import (
	"github.com/gofunct/mamba/docker/commands"
	"github.com/spf13/cobra"
)

var (
	list bool
)

func init() {
	ImageCmd = imageCmd()
}

var ImageCmd *cobra.Command

func imageCmd() *cobra.Command {
	var imgcmd = &cobra.Command{
		Use:   "image",
		Short: "Commands for interacting with a docker images",
	}
	imgcmd.Flags().BoolVarP(&list, "list", "l", false, "list docker images")

	switch {
	case list == true:
		imgcmd.Run = commands.ListImages()
	default:
		imgcmd.Run = commands.DefaultFunc()
	}

	return imgcmd
}
