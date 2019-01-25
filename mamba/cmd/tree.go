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

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
)

func init() {

}

var fs = afero.NewOsFs()
var base = afero.NewBasePathFs(fs, ".")

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use: "tree",
	Run: func(cmd *cobra.Command, args []string) {

		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/", DefaultRouter)
		router.HandleFunc("/files", FilesDirRouter)
		log.Info("Starting contain server on:", "0.0.0.0:10000")
		log.Infof("\nRoutes:\n%s\n%s", "/", "/files")
		log.Info("Press Crtl-C to shutdown...")
		if err := http.ListenAndServe("0.0.0.0:10000", router); err != nil {
			log.Fatalf("%s\n", err)
		}

	},
}

//different router mapping
func DefaultRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome! This is my restful api program! \n")
}

func FilesDirRouter(w http.ResponseWriter, r *http.Request) {

	//call ParseDirTree
	tree, err := ParseDirTree(os.Getenv(("PWD")))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error traversing the filesystem: %v\n", err)
		os.Exit(3)
	} else {

		if err := json.NewEncoder(w).Encode(tree); err != nil {
			panic(err)
		}

	}
}

//def struct DirTree
type DirTree struct {
	IsDir    bool       `json:"IsDir"`
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Children []*DirTree `json:"children"`
}

// build the directory tree
func ParseDirTree(root string) (result *DirTree, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*DirTree)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		parents[path] = &DirTree{
			IsDir:    info.IsDir(),
			Name:     info.Name(),
			Path:     path,
			Children: make([]*DirTree, 0),
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		log.Warnln("filepath failure", err.Error())
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists {
			result = node
		} else {
			parent.Children = append(parent.Children, node)
		}
	}
	return
}

// convert struct to json
func (parsed *DirTree) ToJson() string {
	j, err := json.Marshal(parsed)
	if err != nil {
		log.Infoln("JSON ERROR: " + err.Error())
		return "JSON ERROR: " + err.Error()
	}
	return string(j)
}
