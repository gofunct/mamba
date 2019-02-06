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
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

// Licenses contains all possible licenses a user can choose from.
var Licenses = make(map[string]License)

type License struct {
	Name            string   // The type of license in use
	PossibleMatches []string // Similar names to guess
	Text            string   // License text data
	Header          string   // License header for source files
}

func init() {
	initApache2()
}

func getLicense() License {
	// If explicitly flagged, use that.
	if userLicense != "" {
		return findLicense(userLicense)
	}

	// If user wants to have custom license, use that.
	if viper.IsSet("license.header") || viper.IsSet("license.text") {
		return License{Header: viper.GetString("license.header"),
			Text: viper.GetString("license.text")}
	}

	// If user wants to have built-in license, use that.
	if viper.IsSet("license") {
		return findLicense(viper.GetString("license"))
	}

	// If user didn't set any license, use Apache 2.0 by default.
	return Licenses["apache"]
}

func copyrightLine() string {
	author := os.Getenv("USER")
	year := time.Now().Format("2006")
	return "Copyright © " + year + " " + author
}

// error will be printed.
func findLicense(name string) License {
	found := matchLicense(name)
	if found == "" {
		er("unknown license: " + name)
	}
	return Licenses[found]
}

func matchLicense(name string) string {
	if name == "" {
		return ""
	}

	for key, lic := range Licenses {
		for _, match := range lic.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}
