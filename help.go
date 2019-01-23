package mamba

// HelpFunc returns either the function set by SetHelpFunc for this command
// or a parent, or it returns a function with default help behavior.
func (c *Command) HelpFunc() func(*Command, []string) {
	if c.helpFunc != nil {
		return c.helpFunc
	}
	if c.hasParent() {
		return c.Parent().HelpFunc()
	}
	return func(c *Command, a []string) {
		c.mergePersistentFlags()
		err := tmpl(c.OutOrStdout(), c.HelpTemplate(), c)
		if err != nil {
			c.Println(err)
		}
	}
}

// Help puts out the help for the command.
// Used when a user calls help [command].
// Can be defined by user by overriding HelpFunc.
func (c *Command) Help() error {
	c.HelpFunc()(c, []string{})
	return nil
}

// SetHelpFunc sets help function. Can be defined by Application.
func (c *Command) SetHelpFunc(f func(*Command, []string)) {
	c.helpFunc = f
}

// SetHelpCommand sets help command.
func (c *Command) SetHelpCommand(cmd *Command) {
	c.helpCommand = cmd
}

// SetHelpTemplate sets help template to be used. Application can use it to set custom template.
func (c *Command) SetHelpTemplate(s string) {
	c.helpTemplate = s
}

// HelpTemplate return help template for the command.
func (c *Command) HelpTemplate() string {
	if c.helpTemplate != "" {
		return c.helpTemplate
	}

	if c.hasParent() {
		return c.parent.HelpTemplate()
	}
	return `{{with (or .Info)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

func (c *Command) initDefaultHelpCmd() {
	if !c.hasSubCommands() {
		return
	}

	if c.helpCommand == nil {
		c.helpCommand = &Command{
			Use: "help [command]",
			Info: `Help provides help for any command in the application.
Simply type ` + c.Name() + ` help [path to command] for full details.`,

			Run: func(c *Command, args []string) {
				cmd, _, e := c.Root().find(args)
				if cmd == nil || e != nil {
					c.Printf("Unknown help topic %#q\n", args)
					c.Root().Usage()
				} else {
					cmd.initDefaultHelpFlag() // make possible 'help' flag to be shown
					cmd.Help()
				}
			},
		}
	}
	c.RemoveCommand(c.helpCommand)
	c.AddCommand(c.helpCommand)
}

func (c *Command) initDefaultHelpFlag() {
	c.mergePersistentFlags()
	if c.Flags().Lookup("help") == nil {
		usage := "help for "
		if c.Name() == "" {
			usage += "this command"
		} else {
			usage += c.Name()
		}
		c.Flags().BoolP("help", "h", false, usage)
	}
}
