package gcloud

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type gcloud struct {
	project string
}

type key struct {
	PrivateKeyID string `json:"private_key_id"`
}

type service struct{ Status *status }
type status struct{ LoadBalancer loadBalancer }
type loadBalancer struct{ Ingress []ingress }
type ingress struct{ IP string }

type tfItem struct {
	Sensitive bool
	Type      string
	Value     string
}
type state struct {
	Project          tfItem
	ClusterName      tfItem `json:"cluster_name"`
	ClusterZone      tfItem `json:"cluster_zone"`
	Bucket           tfItem
	DatabaseInstance tfItem `json:"database_instance"`
	DatabaseRegion   tfItem `json:"database_region"`
	MotdVarConfig    tfItem `json:"motd_var_config"`
	MotdVarName      tfItem `json:"motd_var_name"`
}

func run(args ...string) (stdout string, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return strings.TrimSpace(string(stdoutb)), nil
}
func (gcp *gcloud) args(args ...string) []string {
	return append([]string{"gcloud", "--quiet", "--project", gcp.project}, args...)
}

func runString(args ...string) (stdout string, err error) {
	stdoutb, err := runBytes(args...)
	return strings.TrimSpace(string(stdoutb)), err
}

func (gcp *gcloud) cmd(args ...string) *exec.Cmd {
	args = append([]string{"--quiet", "--project", gcp.project}, args...)
	cmd := exec.Command("gcloud", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Stderr = os.Stderr
	return cmd
}

func runBytes(args ...string) (stdout []byte, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return stdoutb, nil
}
