package cli

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// Run parses the gorram command line args and runs the gorram command.
func Run(cmd *cobra.Command) int {
	env := OSEnv{
		Args:   make([]string, len(os.Args)),
		Stderr: cmd.OutOrStdout(),
		Stdout: cmd.OutOrStdout(),
		Stdin:  os.Stdin,
		Env:    getenv(os.Environ()),
	}
	copy(env.Args, os.Args)
	return ParseAndRun(env)
}

func getenv(env []string) map[string]string {
	ret := make(map[string]string, len(env))
	for _, s := range env {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			panic("invalid environment variable set: " + s)
		}
		ret[parts[0]] = parts[1]
	}
	return ret
}
