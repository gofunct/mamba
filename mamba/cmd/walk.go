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
	"github.com/gofunct/mamba/pkg/generator"
	"github.com/shiyanhui/hero"
)

func init() {
	walkCmd.AddCommand(protoGenCmd, goGoCmd, grpcCmd, htmlCmd)
	protoGenCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "path to input directory")
	protoGenCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "path to output directory")
	goGoCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "path to input directory")
	goGoCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "path to output directory")
	htmlCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "path to input directory")
	htmlCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "path to output directory")
	htmlCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "", "package name")
}

// protocGenCmd represents the protocGen command
var protoGenCmd = &mamba.Command{
	Use:  "protocGen",
	Info: "Compile templates as a protoc plugin",
	Run: func(cmd *mamba.Command, args []string) {
		var g = generator.NewGenerator()
		g.Generate(in, out)
	},
}

// protocCmd represents the protoc command
var goGoCmd = &mamba.Command{
	Use:  "gogo",
	Info: "Compile templates as a protoc plugin",
	Run: func(cmd *mamba.Command, args []string) {
		if err := cmd.WalkGoGoProto(in); err != nil {
			fmt.Printf("%+v", err)
		}
	},
}

// protocCmd represents the protoc command
var grpcCmd = &mamba.Command{
	Use:  "grpc",
	Info: "Compile templates as a protoc plugin",
	Run: func(cmd *mamba.Command, args []string) {
		if err := cmd.WalkGoGoProto(in); err != nil {
			fmt.Printf("%+v", err)
		}
	},
}

// protocCmd represents the protoc command
var walkCmd = &mamba.Command{
	Use:  "walk",
	Info: "walk a filepath with a given function and file extension",
}

var htmlCmd = &mamba.Command{
	Use:  "html",
	Info: "Generate html files",
	Run: func(cmd *mamba.Command, args []string) {
		generate(in, out, pkg)
	},
}

func generate(source, dest, pkg string) {
	hero.Generate(source, dest, pkg)
}
