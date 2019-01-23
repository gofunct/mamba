package mamba

import "bytes"

// UsageFunc returns either the function set by SetUsageFunc for this command
// or a parent, or it returns a default usage function.
func (c *Command) UsageFunc() (f func(*Command) error) {
	if c.UsageF != nil {
		return c.UsageF
	}
	if c.hasParent() {
		return c.Parent().UsageFunc()
	}
	return func(c *Command) error {
		c.mergePersistentFlags()
		err := tmpl(c.OutOrStderr(), c.UsageTemplate(), c)
		if err != nil {
			c.Println(err)
		}
		return err
	}
}

// Usage puts out the usage for the command.
// Used when a user provides invalid input.
// Can be defined by user by overriding UsageFunc.
func (c *Command) Usage() error {
	return c.UsageFunc()(c)
}

// UsageTemplate returns usage template for the command.
func (c *Command) UsageTemplate() string {
	if c.UsageTmpl != "" {
		return c.UsageTmpl
	}

	if c.hasParent() {
		return c.parent.UsageTemplate()
	}
	return `
Usage:

{{.UseLine}}
{{.CommandPath}} [command]


Available Commands:{{range .Commands}}
{{.CommandPath}}

{{.Name}} 

{{.Info}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}

`
}

// UsageString return usage string.
func (c *Command) UsageString() string {
	tmpOutput := c.output
	bb := new(bytes.Buffer)
	c.SetOutput(bb)
	c.Usage()
	c.output = tmpOutput
	return bb.String()
}

// SetUsageFunc sets usage function. Usage can be defined by application.
func (c *Command) SetUsageFunc(f func(*Command) error) {
	c.UsageF = f
}

// SetUsageTemplate sets usage template. Can be defined by Application.
func (c *Command) SetUsageTemplate(s string) {
	c.UsageTmpl = s
}
