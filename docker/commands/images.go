package commands

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gofunct/mamba/docker/client"
	"github.com/gofunct/mamba/function"
	"github.com/gofunct/mamba/logging"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"time"
)

func ListImages() function.RunFunc {
	var err error
	return func(cmd *cobra.Command, args []string) {
		if err != nil {
			logging.L.Fatalln(err.Error())
		}
		imgs, err := client.Client.ListImages(docker.ListImagesOptions{All: false})
		if err != nil {
			logging.L.Fatalln(err.Error())
		}
		for _, img := range imgs {
			fmt.Println("ID: ", img.ID)
			fmt.Println("RepoTags: ", img.RepoTags)
			fmt.Println("Created: ", img.Created)
			fmt.Println("Size: ", img.Size)
			fmt.Println("VirtualSize: ", img.VirtualSize)
			fmt.Println("ParentId: ", img.ParentID)
		}
	}
}
func BuildImage() function.RunFunc {

	return func(cmd *cobra.Command, args []string) {
		t := time.Now()
		inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
		tr := tar.NewWriter(inputbuf)
		err := tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
		if err != nil {
			logging.L.Fatalln("failed to write header to dockerfile", errors.WithStack(err).Error())
		}
		i, err := tr.Write([]byte("FROM base\n"))
		if err != nil {
			logging.L.Fatalln("failed to write header to dockerfile", errors.WithStack(err).Error(), i)
		}

		err = tr.Close()
		if err != nil {
			logging.L.Fatalln("failed to write header to dockerfile", errors.WithStack(err).Error())
		}
		opts := docker.BuildImageOptions{
			Name:         "test",
			InputStream:  inputbuf,
			OutputStream: outputbuf,
		}
		if err := client.Client.BuildImage(opts); err != nil {
			logging.L.Fatalln("failed to build image", errors.WithStack(err).Error())
		}
	}
}
