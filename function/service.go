package function

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

type RunFunc func(cmd *cobra.Command, args []string)

func RunString(args ...string) (stdout string, err error) {
	stdoutb, err := RunBytes(args...)
	return strings.TrimSpace(string(stdoutb)), err
}

func RunBytes(args ...string) (stdout []byte, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return stdoutb, nil
}

func ValidateString(args ...string) {

	s, err := RunString(args...)
	if s != "" {
		fmt.Println(s)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ValidateBytes(args ...string) {

	s, err := RunBytes(args...)
	if s != nil {
		fmt.Println(string(s))
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
