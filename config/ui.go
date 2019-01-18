package config

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/prometheus/common/log"
	"strings"
)

func (v *config) Debug() {
	color.Red("Overriden Variables:\n")
	for k, v := range v.override {
		color.Blue(indentcol, k)
		color.Black(indent,v)
	}
	color.Red("Flag Variables:\n")
	for k, v := range v.pflags {
		color.Blue(indentcol, k)
		color.Black(indent,v)
	}
	color.Red("Environmental Variables:\n")
	for k, v := range v.env {
		color.Blue(indentcol, k)
		color.Black(indent,v)
	}
	color.Red("Config Variables:\n")
	for k, v := range v.config {
		color.Blue(indentcol, k)
		color.Black(indent,v)
	}
	color.Red("Default Variables:\n")
	for k, v := range v.defaults {
		color.Blue(indentcol, k)
		color.Black(indent,v)
	}
}

// Search all configPaths for any config file.
// Returns the first path that exists (and is a config file).
func (v *config) findconfigFile() (string, error) {
	log.Debug("event", "finding config file...")

	file := v.searchInPath(v.Path)
	if file != "" {
		return file, nil
	}
	return "", configFileNotFoundError{v.Name, fmt.Sprintf("%s", v.Path)}
}

func (v *config) AllSettings() map[string]interface{} {

	m := map[string]interface{}{}
	// start from the list of keys, and construct the map one value at a time
	for _, k := range v.AllKeys() {
		value := v.Get(k)
		if value == nil {
			// should not happen, since AllKeys() returns only keys holding a value,
			// check just in case anything changes
			continue
		}
		path := strings.Split(k, v.keyDelim)
		lastKey := strings.ToLower(path[len(path)-1])
		deepestMap := deepSearch(m, path[0:len(path)-1])
		// set innermost value
		deepestMap[lastKey] = value
	}
	return m
}

var indent = "%s"
var indentcol = "%s:"
