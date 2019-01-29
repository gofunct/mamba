package function

import (
	"bytes"
	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/gofunct/mamba/runtime"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func RunDgraph() runtime.Function {
	return func(cmd *cobra.Command, args []string) {
		var b []byte
		zap.Debug("pulling dgraph docker image...", "args", `"docker", "pull", "dgraph/dgraph"`)
		pull, err := runtime.RunBytes("docker", "pull", "dgraph/dgraph")
		b = append(b, pull...)
		if err != nil {
			zap.LogF("failed to execute:", err)
		}
		zap.Debug("making directory for data at ~/dgraph...", "args", `"mkdir", "-p", "~/dgraph"`)
		mkdir, err := runtime.RunBytes("mkdir", "-p", "~/dgraph")
		b = append(b, mkdir...)
		if err != nil {
			zap.LogF("failed to execute:", err)
		}
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
		zap.Debug("executing dgraph docker container...", "args", strings.Join(dockerArgs, ""))
		run, err := runtime.RunBytes(dockerArgs...)
		b = append(b, run...)
		if err != nil {
			zap.LogF("failed to execute:", err)
		}
		zap.Debug("executing dgraph docker container...", "args", `"docker", "exec", "-it", "dgraph dgraph alpha", "--lru_mb", "2048", "--zero", "localhost:5080"`)
		exec, err := runtime.RunBytes("docker", "exec", "-it", "dgraph dgraph alpha", "--lru_mb", "2048", "--zero", "localhost:5080")
		b = append(b, exec...)
		if err != nil {
			zap.LogF("failed to execute:", err)
		}
		zap.Debug("executing dgraph docker container...", "args", `"docker", "exec", "-it", "dgraph dgraph-ratel"`)
		ratel, err := runtime.RunBytes("docker", "exec", "-it", "dgraph dgraph-ratel")
		b = append(b, ratel...)
		if err != nil {
			zap.LogF("failed to execute:", err)
		}

		buf := bytes.NewBuffer(b)
		_, err = buf.WriteTo(os.Stdout)
		if err != nil {
			zap.LogF("failed to execute:", err)
		}
	}
}
