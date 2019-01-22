package mamba

import "fmt"

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
