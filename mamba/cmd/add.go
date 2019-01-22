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
	"github.com/gofunct/mamba"
	"os"
	"path/filepath"
	"unicode"
)

func init() {
	addCmd.Flags().StringVarP(&packageName, "package", "t", "", "target package name (e.g. github.com/spf13/hugo)")
	addCmd.Flags().StringVarP(&parentName, "parent", "p", "rootCmd", "variable name of parent command for this command")
}

var packageName, parentName string

var addCmd = &mamba.Command{
	Use:     "add [command name]",
	Aliases: []string{"command"},
	Info: `Add (mamba add) will create a new command, with a license and
the appropriate structure for a Mamba-based CLI application,
and register it to its parent (default rootCmd).

If you want your command to be public, pass in the command name
with an initial uppercase letter.

Example: mamba add server -> resulting in a new cmd/server.go`,

	Run: func(cmd *mamba.Command, args []string) {
		if len(args) < 1 {
			er("add needs a name for the command")
		}

		var project *Project
		if packageName != "" {
			project = NewProject(packageName)
		} else {
			wd, err := os.Getwd()
			if err != nil {
				er(err)
			}
			project = NewProjectFromPath(wd)
		}

		cmdName := validateCmdName(args[0])
		cmdPath := filepath.Join(project.CmdPath(), cmdName+".go")
		createCmdFile(project.License(), cmdPath, cmdName)

		fmt.Fprintln(cmd.OutOrStdout(), cmdName, "created at", cmdPath)
	},
}

func validateCmdName(source string) string {
	i := 0
	l := len(source)
	// The output is initialized on demand, then first dash or underscore
	// occurs.
	var output string

	for i < l {
		if source[i] == '-' || source[i] == '_' {
			if output == "" {
				output = source[:i]
			}

			// If it's last rune and it's dash or underscore,
			// don't add it output and break the loop.
			if i == l-1 {
				break
			}

			// If next character is dash or underscore,
			// just skip the current character.
			if source[i+1] == '-' || source[i+1] == '_' {
				i++
				continue
			}

			// If the current character is dash or underscore,
			// upper next letter and add to output.
			output += string(unicode.ToUpper(rune(source[i+1])))
			// We know, what source[i] is dash or underscore and source[i+1] is
			// uppered character, so make i = i+2.
			i += 2
			continue
		}

		// If the current character isn't dash or underscore,
		// just add it.
		if output != "" {
			output += string(source[i])
		}
		i++
	}

	if output == "" {
		return source // source is initially valid name.
	}
	return output
}

func createCmdFile(license License, path, cmdName string) {
	template := `{{comment .copyright}}
{{if .license}}{{comment .license}}{{end}}

package {{.cmdPackage}}

import (
	"fmt"
	"github.com/spf13/mamba"
)

// {{.cmdName}}Cmd represents the {{.cmdName}} command
var {{.cmdName}}Cmd = &mamba.Command{
	Use:   "{{.cmdName}}",
	Info: "A brief description of your command",
	// Add all config requirements to Valid Args
	ValidArgs: []string{},
	Run: func(cmd *mamba.Command, args []string) {
		
		fmt.Println("{{.cmdName}} valid arguments are: ", cmd.ValidArgs)
	},
}

func init() {
	{{.parentName}}.AddCommand({{.cmdName}}Cmd)
}
`

	data := make(map[string]interface{})
	data["copyright"] = copyrightLine()
	data["license"] = license.Header
	data["cmdPackage"] = filepath.Base(filepath.Dir(path)) // last dir of path
	data["parentName"] = parentName
	data["cmdName"] = cmdName

	cmdScript, err := executeTemplate(template, data)
	if err != nil {
		er(err)
	}
	err = writeStringToFile(path, cmdScript)
	if err != nil {
		er(err)
	}
}
