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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:     "init [name]",
	Aliases: []string{"initialize", "initialise", "create"},
	Long: `Initialize (mamba init) will create a new application, with a license
and the appropriate structure for a Mamba-based CLI application.

  * If a name is provided, it will be created in the current directory;
  * If no name is provided, the current directory will be assumed;
  * If a relative path is provided, it will be created inside $GOPATH
    (e.g. github.com/spf13/hugo);
  * If an absolute path is provided, it will be created;
  * If the directory already exists but is empty, it will be used.

Init will not use an existing directory with contents.`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			er(err)
		}

		var project *Project
		if len(args) == 0 {
			project = NewProjectFromPath(wd)
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			if filepath.IsAbs(arg) {
				project = NewProjectFromPath(arg)
			} else {
				project = NewProject(arg)
			}
		} else {
			er("please provide only one argument")
		}

		initializeProject(project)
		OsExec("go", "mod", "init")
		OsExec("go", "mod", "vendor")
		OsExec("go", "fmt", "./...")
		OsExec("go", "install")
		fmt.Fprintln(cmd.OutOrStdout(), `Your Mamba application is ready at
`+project.AbsPath()+`

Give it a try by going there and running `+"`go run main.go`."+`
Add commands to it by running `+"`mamba add [cmdname]`.")
	},
}

func initializeProject(project *Project) {
	if !exists(project.AbsPath()) { // If path doesn't yet exist, create it
		err := os.MkdirAll(project.AbsPath(), os.ModePerm)
		if err != nil {
			er(err)
		}
	} else if !isEmpty(project.AbsPath()) { // If path exists and is not empty don't use it
		er("Mamba will not create a new project in a non empty directory: " + project.AbsPath())
	}

	// We have a directory and it's empty. Time to initialize it.
	createLicenseFile(project.License(), project.AbsPath())
	createMainFile(project)
	createRootCmdFile(project)
}

func createLicenseFile(license License, path string) {
	data := make(map[string]interface{})
	data["copyright"] = copyrightLine()

	// Generate license template from text and data.
	text, err := executeTemplate(license.Text, data)
	if err != nil {
		er(err)
	}

	// Write license text to LICENSE file.
	err = writeStringToFile(filepath.Join(path, "LICENSE"), text)
	if err != nil {
		er(err)
	}
}

func createMainFile(project *Project) {
	mainTemplate := `{{ comment .copyright }}
{{if .license}}{{ comment .license }}{{end}}

package main

import "{{ .importpath }}"

func main() {
	cmd.Execute()
}
`
	data := make(map[string]interface{})
	data["copyright"] = copyrightLine()
	data["license"] = project.License().Header
	data["importpath"] = path.Join(project.Name(), filepath.Base(project.CmdPath()))

	mainScript, err := executeTemplate(mainTemplate, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(project.AbsPath(), "main.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func createRootCmdFile(project *Project) {
	template := `{{comment .copyright}}
{{if .license}}{{comment .license}}{{end}}

package cmd

import (
	"github.com/gofunct/mamba"
	"github.com/pkg/errors"
	"context"
	"fmt"
	"os"
	"net/http"
)

func init() {
	ctx = context.TODO()
}

var (
	ctx context.Context
)

var root = &mamba.Command{
	Version:      "v0.1.1",
	Dependencies: nil,
	PreRun: func(svc *mamba.Command, ctx context.Context) {
		fmt.Println("Welcome "+os.Getenv("USER")+"!")
	},
	Login: func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "this is where your users will login")
	},
	Home: func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "this is where your web app will be located")

	},
	FAQ: func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "this is where your users will go for help")
	},
}

func Execute() {
	if err := root.Execute(ctx); err != nil {
		panic(errors.WithStack(err))
	}
}
`

	data := make(map[string]interface{})
	data["copyright"] = copyrightLine()
	data["license"] = project.License().Header
	data["appName"] = path.Base(project.Name())

	rootCmdScript, err := executeTemplate(template, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(project.CmdPath(), "root.go"), rootCmdScript)
	if err != nil {
		er(err)
	}

}
