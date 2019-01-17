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
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	RootCmd.AddCommand(dbCmd)
}

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Db()
	},
}

func Db() {
	log.SetFlags(0)
	log.SetPrefix("gcp/provision_db: ")
	project := flag.String("project", "", "GCP project ID")
	serviceAccount := flag.String("service_account", "", "name of service account in GCP project")
	instance := flag.String("instance", "", "database instance name")
	database := flag.String("database", "", "name of database to initialize")
	password := flag.String("password", "", "root password for the database")
	schema := flag.String("schema", "", "path to .sql file defining the database schema")
	flag.Parse()
	missing := false
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			log.Printf("Required flag -%s is not set.", f.Name)
			missing = true
		}
	})
	if missing {
		os.Exit(64)
	}
	if err := provisionDB(*project, *serviceAccount, *instance, *database, *password, *schema); err != nil {
		log.Fatal(err)
	}
}

func provisionDB(projectID, serviceAccount, dbInstance, DbName, DbPassword, schemaPath string) error {
	log.Printf("Downloading Docker images...")
	const mySQLImage = "mysql:5.6"
	cloudSQLProxyImage := "gcr.io/cloudsql-docker/gce-proxy:1.11"
	images := []string{mySQLImage, cloudSQLProxyImage}
	for _, img := range images {
		if _, err := run("docker", "pull", img); err != nil {
			return err
		}
	}

	log.Printf("Getting connection string from database metadata...")
	gcp := &gcloud{projectID}
	dbConnStr, err := run(gcp.args("sql", "instances", "describe", "--format", "value(connectionName)", dbInstance)...)
	if err != nil {
		return fmt.Errorf("getting connection string: %v", err)
	}

	// Create a temporary directory to hold the service account key.
	// We resolve all symlinks to avoid Docker on Mac issues, see
	// https://github.com/google/go-cloud/issues/110.
	serviceAccountVolDir, err := ioutil.TempDir("", "guestbook-service-acct")
	if err != nil {
		return fmt.Errorf("creating temp dir to hold service account key: %v", err)
	}
	serviceAccountVolDir, err = filepath.EvalSymlinks(serviceAccountVolDir)
	if err != nil {
		return fmt.Errorf("evaluating any symlinks: %v", err)
	}
	defer os.RemoveAll(serviceAccountVolDir)
	log.Printf("Created %v", serviceAccountVolDir)

	// Furnish a new service account key.
	if _, err := run(gcp.args("iam", "service-accounts", "keys", "create", "--iam-account="+serviceAccount, serviceAccountVolDir+"/key.json")...); err != nil {
		return fmt.Errorf("creating new service account key: %v", err)
	}
	keyJSONb, err := ioutil.ReadFile(filepath.Join(serviceAccountVolDir, "key.json"))
	if err != nil {
		return fmt.Errorf("reading key.json file: %v", err)
	}
	var k key
	if err := json.Unmarshal(keyJSONb, &k); err != nil {
		return fmt.Errorf("parsing key.json: %v", err)
	}
	serviceAccountKeyID := k.PrivateKeyID
	defer func() {
		if _, err := run(gcp.args("iam", "service-accounts", "keys", "delete", "--iam-account", serviceAccount, serviceAccountKeyID)...); err != nil {
			log.Printf("deleting service account key: %v", err)
		}
	}()
	log.Printf("Created service account key %s", serviceAccountKeyID)

	log.Printf("Starting Cloud SQL proxy...")
	proxyContainerID, err := run("docker", "run", "--detach", "--rm", "--volume", serviceAccountVolDir+":/creds", "--publish", "3306", cloudSQLProxyImage, "/cloud_sql_proxy", "-instances", dbConnStr+"=tcp:0.0.0.0:3306", "-credential_file=/creds/key.json")
	if err != nil {
		return err
	}
	defer func() {
		if _, err := run("docker", "kill", proxyContainerID); err != nil {
			log.Printf("failed to kill docker container for proxy: %v", err)
		}
	}()

	log.Print("Sending schema to database...")
	mySQLCmd := fmt.Sprintf(`mysql --wait -h"$PROXY_PORT_3306_TCP_ADDR" -P"$PROXY_PORT_3306_TCP_PORT" -uroot -p'%s' '%s'`, DbPassword, DbName)
	connect := exec.Command("docker", "run", "--rm", "--interactive", "--link", proxyContainerID+":proxy", mySQLImage, "sh", "-c", mySQLCmd)
	schema, err := os.Open(schemaPath)
	if err != nil {
		return err
	}
	defer schema.Close()
	connect.Stdin = schema
	connect.Stderr = os.Stderr
	if err := connect.Run(); err != nil {
		return fmt.Errorf("running %v: %v", connect.Args, err)
	}

	return nil
}
