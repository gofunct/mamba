package logging

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	L = NewLogCtx(logrus.New())
	L.SetOutput(os.Stdout)
	L.SetFormatter(&logrus.JSONFormatter{})
}

var L *contextLog

func NewLogCtx(logger *logrus.Logger) *contextLog {
	return &contextLog{
		lgr:     lgr{logger},
		Context: NewLogrusLogger(logger),
	}
}

type contextLog struct {
	lgr
	Context log.Logger
}

// NewLogrusLogger returns a go-kit log.lgr that sends log events to a Logrus logger.
func NewLogrusLogger(logger *logrus.Logger) log.Logger {
	return &lgr{
		logger,
	}
}

func (l lgr) Log(keyvals ...interface{}) error {
	fields := logrus.Fields{}
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			fields[fmt.Sprint(keyvals[i])] = keyvals[i+1]
		} else {
			fields[fmt.Sprint(keyvals[i])] = errMissingValue
		}
	}
	l.WithFields(fields).Info()
	return nil
}

type lgr struct {
	*logrus.Logger
}

var (
	DebugLevel = func() {
		L.SetLevel(logrus.DebugLevel)
	}
	WarnLevel = func() {
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
			L.Log("cmd", cmd.Name())
		case warnEnabled:
			WarnLevel()
			L.Log("cmd", cmd.Name())
		default:
			DebugLevel()
			L.Log("cmd", cmd.Name())
		}
	})
}

var errMissingValue = errors.New("(MISSING)")
