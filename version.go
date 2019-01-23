package mamba

func (c *Command) initDefaultVersionFlag() {
	if c.Version == "" {
		return
	}

	c.mergePersistentFlags()
	if c.Flags().Lookup("version") == nil {
		usage := "version for "
		if c.Name() == "" {
			usage += "this command"
		} else {
			usage += c.Name()
		}
		c.Flags().Bool("version", false, usage)
	}
}

// VersionTemplate return version template for the command.
func (c *Command) VersionTemplate() string {
	if c.versionTemplate != "" {
		return c.versionTemplate
	}

	if c.hasParent() {
		return c.parent.VersionTemplate()
	}
	return `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`
}

// SetVersionTemplate sets version template to be used. Application can use it to set custom template.
func (c *Command) SetVersionTemplate(s string) {
	c.versionTemplate = s
}
