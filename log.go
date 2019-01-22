package mamba

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
