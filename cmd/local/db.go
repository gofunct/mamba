// Copyright © 2019 Coleman Word <coleman.word@gofunct.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package local

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/pipe.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(dbCmd)
}

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Service()
	},
}

func Service() {
	guestbookDir := flag.String("guestbook_dir", ".", "directory containing guestbook sample source code")
	flag.Parse()
	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "usage: localdb [c] container_name\n")
		os.Exit(1)
	}
	log.SetPrefix("localdb: ")
	log.SetFlags(0)
	if err := runLocalDB(flag.Arg(0), *guestbookDir); err != nil {
		log.Fatal(err)
	}
}

func runLocalDB(containerName, guestbookDir string) error {
	image := "mysql:5.6"

	log.Printf("Starting container running MySQL")
	dockerArgs := []string{"run", "--rm"}
	if containerName != "" {
		dockerArgs = append(dockerArgs, "--name", containerName)
	}
	dockerArgs = append(dockerArgs,
		"--env", "MYSQL_DATABASE=guestbook",
		"--env", "MYSQL_ROOT_PASSWORD=password",
		"--detach",
		"--publish", "3306:3306",
		image)
	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("running %v: %v: %s", cmd.Args, err, out)
	}
	containerID := strings.TrimSpace(string(out))
	defer func() {
		log.Printf("killing %s", containerID)
		stop := exec.Command("docker", "kill", containerID)
		stop.Stderr = os.Stderr
		if err := stop.Run(); err != nil {
			log.Printf("failed to kill db container: %v", err)
		}
	}()

	// Stop the container on Ctrl-C.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		// TODO(ijt): Handle SIGTERM.
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	nap := 10 * time.Second
	log.Printf("Waiting %v for database to come up", nap)
	select {
	case <-time.After(nap):
		// ok
	case <-ctx.Done():
		return errors.New("interrupted while napping")
	}

	log.Printf("Initializing database schema and users")
	schema, err := ioutil.ReadFile(filepath.Join(guestbookDir, "schema.sql"))
	if err != nil {
		return fmt.Errorf("reading schema: %v", err)
	}
	roles, err := ioutil.ReadFile(filepath.Join(guestbookDir, "roles.sql"))
	if err != nil {
		return fmt.Errorf("reading roles: %v", err)
	}
	tooMany := 10
	var i int
	for i = 0; i < tooMany; i++ {
		mySQL := `mysql -h"${MYSQL_PORT_3306_TCP_ADDR?}" -P"${MYSQL_PORT_3306_TCP_PORT?}" -uroot -ppassword guestbook`
		p := pipe.Line(
			pipe.Read(strings.NewReader(string(schema)+string(roles))),
			pipe.Exec("docker", "run", "--rm", "--interactive", "--link", containerID+":mysql", image, "sh", "-c", mySQL),
		)
		if _, stderr, err := pipe.DividedOutput(p); err != nil {
			log.Printf("Failed to seed database: %q; retrying", stderr)
			select {
			case <-time.After(time.Second):
				continue
			case <-ctx.Done():
				return errors.New("interrupted while napping in between database seeding attempts")
			}
		}
		break
	}
	if i == tooMany {
		return fmt.Errorf("gave up after %d tries to seed database", i)
	}

	log.Printf("Database running at localhost:3306")
	attach := exec.CommandContext(ctx, "docker", "attach", containerID)
	attach.Stdout = os.Stdout
	attach.Stderr = os.Stderr
	if err := attach.Run(); err != nil {
		return fmt.Errorf("running %v: %q", attach.Args, err)
	}

	return nil
}
