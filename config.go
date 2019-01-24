package mamba

import (
	"fmt"
	"github.com/gofunct/mamba/pkg/input"
	"github.com/gofunct/mamba/pkg/logging"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func (m *Command) Set(s string, k interface{}) {
	viper.SetDefault(s, k)
}

func (m *Command) GetString(s string) string {

	if res := viper.GetString(s); res != "" {
		return res
	}

	if res := os.Getenv(strings.ToUpper(s)); res != "" {
		return res
	}
	q := &input.Query{
		Q: "Please provide a value for: " + s,
		Opts: &input.Options{
			Key:      s,
			Loop:     true,
			Required: true,
		},
	}
	if res, err := query.Ask(q); res != "" && err == nil {
		return res
	} else {
		logging.L.Fatalf("%s failed to retrieve value for: %s\n%s", s, err.Error())
	}
	return ""
}

func (m *Command) SyncRequirements() {
	for _, s := range m.Dependencies {
		if res := viper.GetString(s); res != "" {
			if err := viper.BindEnv(s, res); err != nil {
				logger.Warn("failed to bind env", s, res)
			}
		}

		if res := os.Getenv(strings.ToUpper(s)); res != "" {
			viper.SetDefault(s, res)
			if err := os.Setenv(strings.ToUpper(s), res); err != nil {
				logger.Warn("failed to set env", strings.ToUpper(s), res)
			}
		}
		needed := query.SingleQuery("Please set a value for variable: "+s, &input.Options{
			Key:      s,
			Loop:     true,
			Required: true,
		})
		if res, err := query.Ask(needed); err == nil && res != "" {
			viper.SetDefault(s, res)
			if err := viper.BindEnv(s, res); err != nil {
				logger.Warn("failed to bind env", s, res)
			}
		}
	}
	m.Write()
	logger.Debug("all requirements have been synced")
}

func (m *Command) GetDependencies() map[string]string {
	set := make(map[string]string)
	for _, v := range m.Dependencies {
		set[v] = m.GetString(v)
	}
	return set
}

func (m *Command) DebugQuery() {
	fmt.Println("Config:")
	fmt.Println("Dependencies: ", m.Dependencies)
	for _, v := range query.Queries {
		fmt.Println("Question: ", v.Q)
		fmt.Println("Key: ", v.Opts.Key)
		fmt.Println("Options: ", v.Opts.Options)
		fmt.Println("Required: ", v.Opts.Required)
		fmt.Println("Default: ", v.Opts.Default)
		fmt.Println("Loop: ", v.Opts.Loop)
	}
}

func (m *Command) Write() {
	sel := query.SingleQuery("Write current config to disc?", &input.Options{
		Key:      "config.write",
		Options:  []string{"true", "false"},
		Default:  "false",
		Loop:     true,
		Required: true,
	})
	cfg := query.SingleQuery("Please provide a path to config file", &input.Options{
		Key:      "config.path",
		Options:  []string{"y", "n"},
		Default:  "n",
		Loop:     true,
		Required: true,
	})

	if res, err := query.Select(sel); res == "true" && err == nil {
		if file, err := query.Ask(cfg); err == nil && file != "" {
			if err := viper.WriteConfigAs(file); err != nil {
				logger.Warn("failed to write config", err.Error())
			}
			logger.Debug("Updated config file successfully")
		}
	}
}

func (m *Command) Unmarshal(i interface{}) {
	if err := viper.Unmarshal(i); err != nil {
		logger.Warn("failed to unmarshal", err.Error())
	}
}

func (m *Command) GetMeta() map[string]interface{} {
	if err := viper.MergeInConfig(); err != nil {
		logging.L.Warn("failed to merge in config", err.Error())
	}
	return viper.AllSettings()
}
