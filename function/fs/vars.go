package fs

import (
	"os"
	"os/exec"
	"path/filepath"
)

var ProtoWalkFunc = func(path string, info os.FileInfo, err error) error {
	// skip vendor directory
	if info.IsDir() && info.Name() == "vendor" {
		return filepath.SkipDir
	}
	// find all protobuf files
	if filepath.Ext(path) == ".proto" {
		// args
		args := []string{
			"--go_out=plugins=grpc:.",
			path,
		}
		cmd := exec.Command("protoc", args...)
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
