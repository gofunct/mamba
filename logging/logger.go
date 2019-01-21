package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {

	}

var L *logrus.Logger
// LoggingMode represents a logging configuration specification.
type LoggingMode int

// LoggingMode values
const (
	LoggingNop LoggingMode = iota
	LoggingVerbose
	LoggingDebug
)

var (
	DebugLevel = func() {
		L = logrus.New()
		L.SetOutput(os.Stdout)
		L.SetFormatter(&logrus.JSONFormatter{})
		L.SetLevel(logrus.DebugLevel)
	}
	WarnLevel = func() {
		L = logrus.New()
		L.SetOutput(os.Stdout)
		L.SetFormatter(&logrus.JSONFormatter{})
		L.SetLevel(logrus.WarnLevel)
	}
)

// AddLoggingFlags sets "--debug" and "--verbose" flags to the given *cobra.Command instance.
func AddLoggingFlags(cmd *cobra.Command) {
	var (
		debugEnabled, warnEnabled bool
	)

	cmd.PersistentFlags().BoolVar(
		&debugEnabled,
		"debug",
		false,
		fmt.Sprintf("Debug level output"),
	)
	cmd.PersistentFlags().BoolVar(
		&warnEnabled,
		"warn",
		false,
		fmt.Sprintf("Warn level output"),
	)

	cobra.OnInitialize(func() {
		switch {
		case debugEnabled:
			DebugLevel()
			L.WithField("command", cmd.Name())
		case warnEnabled:
			L.WithField("command", cmd.Name())
			WarnLevel()
		}
	})
}