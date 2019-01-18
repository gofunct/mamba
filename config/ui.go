package config

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-kit/kit/log"
	"os"
	"strings"
)

var (
	red         = color.New(color.FgRed).Add(color.Underline).PrintfFunc()
	redstring   = color.New(color.FgRed).Add(color.Underline).SprintfFunc()
	green       = color.New(color.FgGreen).PrintfFunc()
	greenstring = color.New(color.FgGreen).SprintfFunc()
)

var logger log.Logger

var logkv = func(c *config, k, v string) {
	log.With(logger, "config", c.Name, "path", c.Path)
	if err := logger.Log(greenstring(k), greenstring(v)); err != nil {
		red("global logger failure")
	}
}

func init() {
	{
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "time", log.DefaultTimestampUTC)
	}
}

func (v *config) Debug() {
	green("Aliases:\n%#v\n", v.aliases)
	green("Override:\n%#v\n", v.override)
	green("PFlags:\n%#v\n", v.pflags)
	green("Env:\n%#v\n", v.env)
	green("Key/Value Store:\n%#v\n", v.kvstore)
	green("config:\n%#v\n", v.config)
	green("Defaults:\n%#v\n", v.defaults)
}

// Search all configPaths for any config file.
// Returns the first path that exists (and is a config file).
func (v *config) findconfigFile() (string, error) {
	logkv(v, "event", "finding config file...")

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
