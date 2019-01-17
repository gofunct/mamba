package script

import (
	"context"
	"log"
	"os"
	"os/exec"
)

type ScriptHandler struct{}

func NewScriptHandler() *ScriptHandler {
	return &ScriptHandler{}
}

func (s *ScriptHandler) Exec(ctx context.Context, cmd *Command) (*Output, error) {
	e := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	cmd.Env = os.Environ()
	data, err := e.Output()
	log.Println(string(data))
	return &Output{
		Out: data,
	}, err
}
