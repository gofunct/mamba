package mamba

import (
	"fmt"
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
func (c *Command) AddLoggingFlags() {
	var (
		debugEnabled, warnEnabled bool
	)

	c.PersistentFlags().BoolVar(
		&debugEnabled,
		"debug",
		false,
		fmt.Sprintf("Debug level output"),
	)
	c.PersistentFlags().BoolVar(
		&warnEnabled,
		"warn",
		false,
		fmt.Sprintf("Warn level output"),
	)

	OnInitialize(func() {
		switch {
		case debugEnabled:
			DebugLevel()
			logger.Log("cmd", c.Name())
		case warnEnabled:
			WarnLevel()
			logger.Log("cmd", c.Name())
		default:
			DebugLevel()
			logger.Log("cmd", c.Name())
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

// Print is a convenience method to Print to the defined output, fallback to Stderr if not set.
func (c *Command) Print(i ...interface{}) {
	fmt.Fprint(c.OutOrStderr(), i...)
}

// Println is a convenience method to Println to the defined output, fallback to Stderr if not set.
func (c *Command) Println(i ...interface{}) {
	c.Print(fmt.Sprintln(i...))
}

// Printf is a convenience method to Printf to the defined output, fallback to Stderr if not set.
func (c *Command) Printf(format string, i ...interface{}) {
	c.Print(fmt.Sprintf(format, i...))
}
