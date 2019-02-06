// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofunct/common/pkg/transport/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type DefaultFlags struct {
	Port    string
	Input   string
	Output  string
	Package string
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		r := mux.NewRouter()

		r.Handle("/", http.FileServer(http.Dir("static")))
		r.HandleFunc("/generate", generate)
		addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
		if addr == ":" {
			addr = ":9000"
		}

		fmt.Printf("Listening on %s...\n", addr)
		h := handlers.LoggingHandler(os.Stderr, r)
		h = handlers.CompressHandler(h)
		h = handlers.RecoveryHandler()(h)
		if err := http.ListenAndServe(addr, h); err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
}

func generate(w http.ResponseWriter, r *http.Request) {
	// read input
	decoder := json.NewDecoder(r.Body)
	type Input struct {
		Protobuf string `json:"protobuf"`
		Template string `json:"template"`
	}
	var input Input
	if err := decoder.Decode(&input); err != nil {
		returnError(w, err)
		return
	}

	// create workspace
	dir, err := ioutil.TempDir("", "pggt")
	if err != nil {
		returnError(w, err)
	}
	// clean up
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			log.Printf("error: failed to remove temporary directory: %v", err)
		}
	}()
	if err = ioutil.WriteFile(filepath.Join(dir, "example.proto"), []byte(input.Protobuf), 0644); err != nil {
		returnError(w, err)
		return
	}
	if err = ioutil.WriteFile(filepath.Join(dir, "example.output.tmpl"), []byte(input.Template), 0644); err != nil {
		returnError(w, err)
		return
	}

	// generate
	cmd := exec.Command("protoc", "-I"+dir, "--gotemplate_out=template_dir="+dir+",debug=true:"+dir, filepath.Join(dir, "example.proto")) // #nosec
	out, err := cmd.CombinedOutput()
	if err != nil {
		returnError(w, errors.New(string(out)))
		return
	}

	// read output
	content, err := ioutil.ReadFile(filepath.Join(dir, "example.output")) // #nosec
	if err != nil {
		returnError(w, err)
		return
	}

	returnContent(w, content)
}

func returnContent(w http.ResponseWriter, output interface{}) {
	payload := map[string]interface{}{
		"output": fmt.Sprintf("%s", output),
	}
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func returnError(w http.ResponseWriter, err error) {
	payload := map[string]interface{}{
		"error": fmt.Sprintf("%v", err),
	}
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
