package commands

import (
	"bytes"
	"github.com/gofunct/mamba/function"
	"github.com/gofunct/mamba/manager/logging"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

func RunDgraph() function.RunFunc {
	return func(cmd *cobra.Command, args []string) {
		var b []byte
		logging.L.Debug("pulling dgraph docker image...\n")
		pull, err := function.RunBytes("docker", "pull", "dgraph/dgraph")
		b = append(b, pull...)
		if err != nil {
			logging.L.Fatalln("failed to execute:%s%s", errors.WithStack(err), pull)
		}
		logging.L.Debug("making directory for data at ~/dgraph...\n")
		mkdir, err := function.RunBytes("mkdir", "-p", "~/dgraph")
		b = append(b, mkdir...)
		if err != nil {
			logging.L.Fatalln("failed to execute:%s%s", errors.WithStack(err), mkdir)
		}
		logging.L.Debug("running dgraph docker container...\n")
		dockerArgs := []string{
			"docker", "run", "-i",
			"-p", "6080:6080",
			"-p", "8080:8080",
			"-p", "9080:9080",
			"-p", "8000:8000",
			"-v", "~/dgraph:/dgraph",
			"--name", "dgraph",
			"dgraph/dgraph", "dgraph zero",
		}

		run, err := function.RunBytes(dockerArgs...)
		b = append(b, run...)
		if err != nil {
			logging.L.Fatalln("failed to execute:%s%s", errors.WithStack(err), run)
		}
		logging.L.Debug("executing dgraph docker container...\n")
		exec, err := function.RunBytes("docker", "exec", "-it", "dgraph dgraph alpha", "--lru_mb", "2048", "--zero", "localhost:5080")
		b = append(b, exec...)
		if err != nil {
			logging.L.Fatalln("failed to execute:%s%s", errors.WithStack(err), exec)
		}
		logging.L.Debug("executing dgraph ratel container...\n")
		ratel, err := function.RunBytes("docker", "exec", "-it", "dgraph dgraph-ratel")
		b = append(b, ratel...)
		if err != nil {
			logging.L.Fatalln("failed to execute:%s%s", errors.WithStack(err), ratel)
		}

		buf := bytes.NewBuffer(b)
		_, err = buf.WriteTo(os.Stdout)
		if err != nil {
			logging.L.Fatalln("failed to execute:%s%s", errors.WithStack(err), b)
		}
	}

}
