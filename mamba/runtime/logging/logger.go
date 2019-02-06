package logging

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	L = NewLogCtx(logrus.New())
	L.SetOutput(os.Stdout)
	L.SetFormatter(&logrus.JSONFormatter{})
}

var L *CtxLogger

func NewLogCtx(logger *logrus.Logger) *CtxLogger {
	return &CtxLogger{
		lgr:     lgr{logger},
		Context: NewLogrusLogger(logger),
	}
}

type CtxLogger struct {
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

var errMissingValue = errors.New("(MISSING)")
