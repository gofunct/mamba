package mamba

import (
	"context"
	"github.com/prometheus/common/log"
	osexec "os/exec"
)

func (c *Command) Script(ctx context.Context, args ...string) {
	log.With("cmd", args)
	cmd := osexec.CommandContext(ctx, args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
