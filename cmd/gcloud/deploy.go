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

package gcloud

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(deployCmd)
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Deploy()
	},
}

func Deploy() {
	log.SetFlags(0)
	log.SetPrefix("gcp/deploy: ")
	guestbookDir := flag.String("guestbook_dir", "..", "directory containing the guestbook example")
	tfStatePath := flag.String("tfstate", "terraform.tfstate", "path to terraform state file")
	flag.Parse()
	if err := deploy(*guestbookDir, *tfStatePath); err != nil {
		log.Fatal(err)
	}
}

func deploy(guestbookDir, tfStatePath string) error {
	tfStateb, err := runBytes("terraform", "output", "-state", tfStatePath, "-json")
	if err != nil {
		return err
	}
	var tfState state
	if err := json.Unmarshal(tfStateb, &tfState); err != nil {
		return fmt.Errorf("parsing terraform state JSON: %v", err)
	}
	zone := tfState.ClusterZone.Value
	if zone == "" {
		return fmt.Errorf("empty or missing cluster_zone in %s", tfStatePath)
	}
	tempDir, err := ioutil.TempDir("", "guestbook-k8s-")
	if err != nil {
		return fmt.Errorf("making temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Fill in Kubernetes template parameters.
	proj := strings.Replace(tfState.Project.Value, ":", "/", -1)
	imageName := fmt.Sprintf("gcr.io/%s/guestbook", proj)
	gbyin, err := ioutil.ReadFile(filepath.Join(guestbookDir, "gcp", "guestbook.yaml.in"))
	if err != nil {
		return fmt.Errorf("reading guestbook.yaml.in: %v", err)
	}
	gby := string(gbyin)
	replacements := map[string]string{
		"{{IMAGE}}":             imageName,
		"{{bucket}}":            tfState.Bucket.Value,
		"{{database_instance}}": tfState.DatabaseInstance.Value,
		"{{database_region}}":   tfState.DatabaseRegion.Value,
		"{{motd_var_config}}":   tfState.MotdVarConfig.Value,
		"{{motd_var_name}}":     tfState.MotdVarName.Value,
	}
	for old, new := range replacements {
		gby = strings.Replace(gby, old, new, -1)
	}
	if err := ioutil.WriteFile(filepath.Join(tempDir, "guestbook.yaml"), []byte(gby), 0666); err != nil {
		return fmt.Errorf("writing guestbook.yaml: %v", err)
	}

	// Build Guestbook Docker image.
	log.Printf("Building %s...", imageName)
	build := exec.Command("go", "build", "-o", "gcp/guestbook")
	env := append(build.Env, "GOOS=linux", "GOARCH=amd64")
	env = append(env, os.Environ()...)
	build.Env = env
	absDir, err := filepath.Abs(guestbookDir)
	if err != nil {
		return fmt.Errorf("getting abs path to guestbook dir (%s): %v", guestbookDir, err)
	}
	build.Dir = absDir
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("building guestbook app by running %v: %v", build.Args, err)
	}
	gcp := gcloud{project: tfState.Project.Value}
	cbs := gcp.cmd("builds", "submit", "-t", imageName, filepath.Join(guestbookDir, "gcp"))
	if err := cbs.Run(); err != nil {
		return fmt.Errorf("building container image with %v: %v", cbs.Args, err)
	}

	// Run on Kubernetes.
	log.Printf("Deploying to %s...", tfState.ClusterName.Value)
	getCreds := gcp.cmd("container", "clusters", "get-credentials", "--zone", zone, tfState.ClusterName.Value)
	getCreds.Stderr = os.Stderr
	if err := getCreds.Run(); err != nil {
		return fmt.Errorf("getting credentials with %v: %v", getCreds.Args, err)
	}
	kubeCmds := [][]string{
		{"kubectl", "apply", "-f", filepath.Join(tempDir, "guestbook.yaml")},
		// Force pull the latest image.
		{"kubectl", "scale", "--replicas", "0", "deployment/guestbook"},
		{"kubectl", "scale", "--replicas", "1", "deployment/guestbook"},
	}
	for _, kcmd := range kubeCmds {
		cmd := exec.Command(kcmd[0], kcmd[1:]...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("running %v: %v", cmd.Args, err)
		}
	}

	// Wait for endpoint then print it.
	log.Printf("Waiting for load balancer...")
	for {
		outb, err := runBytes("kubectl", "get", "service", "guestbook", "-o", "json")
		if err != nil {
			return err
		}
		var s service
		if err := json.Unmarshal(outb, &s); err != nil {
			return fmt.Errorf("parsing JSON output: %v", err)
		}
		i := s.Status.LoadBalancer.Ingress
		if len(i) == 0 || i[0].IP == "" {
			dt := time.Second
			log.Printf("No ingress returned in %s. Trying again in %v", outb, dt)
			time.Sleep(dt)
			continue
		}
		endpoint := i[0].IP
		log.Printf("Deployed at http://%s:8080", endpoint)
		break
	}
	return nil
}
