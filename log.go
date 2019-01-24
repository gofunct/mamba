package mamba

import (
	"github.com/gofunct/mamba/pkg/input"
	"github.com/sirupsen/logrus"
)

var (
	DebugLevel = func() {
		logger.SetLevel(logrus.DebugLevel)
	}
	WarnLevel = func() {
		logger.SetLevel(logrus.WarnLevel)
	}
)

// AddLoggingFlags sets "--debug" and "--verbose" flags to the given *cobra.Command instance.
func (c *Command) AddLogging() {

	l, _ := query.Select(&input.Query{
		Q: "what level logging to enable?",
		Opts: &input.Options{
			Key:      "log-level",
			Options:  []string{"warn", "debug"},
			Default:  "debug",
			Loop:     false,
			Required: false,
		},
	})

	OnInitialize(func() {
		switch {
		case l == "debug":
			DebugLevel()
			logger.Log("cmd", c.Version)
		case l == "warn":
			WarnLevel()
			logger.Log("cmd", c.Version)
		default:
			DebugLevel()
			logger.Log("cmd", c.Version)
		}
	})
}

func (m *Command) Warnf(f string, args ...interface{}) {
	logger.Warnf(f, args)
}
func (m *Command) Fatalf(f string, args ...interface{}) {
	logger.Fatalf(f, args)
}

func (m *Command) Debug(args ...interface{}) {
	logger.Debug(args)
}

func (m *Command) Log(args ...interface{}) {
	if err := logger.Log(args); err != nil {
		logger.Warn("failed to log context", err.Error())
	}
}
