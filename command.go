// Copyright © 2019 Coleman Word <coleman.word@gofunct.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package mamba

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

type Command struct {
	// Use is the one-line usage message.
	Use     string
	Version string
	Aliases []string
	// Long is the long message shown in the 'help <this-command>' output.
	Info string
	// Hidden defines, if this command is hidden and should NOT show up in the list of available commands.
	Hidden bool
	// Annotations are key/value pairs that can be used by applications to identify or
	// group commands.
	Env map[string]string
	// ValidArgs is list of all valid non-flag arguments that are accepted in bash completions
	ValidArgs []string

	// Expected arguments
	Args PositionalArgs

	usageFunc             func(*Command) error
	usageTemplate          string
	DisableFlagParsing bool

	PreRun MambaFunc
	// Run: Typically the actual work function. Most commands will only implement this.
	Run MambaFunc
	// PostRun: run after the Run command.
	PostRun MambaFunc

	Router *mux.Router
	// commands is the list of commands supported by this program.
	commands []*Command
	// parent is a parent command for this command.
	parent *Command
	// Max lengths of commands' string lengths for use in padding.
	commandsMaxUseLen         int
	commandsMaxCommandPathLen int
	commandsMaxNameLen        int
	// commandsAreSorted defines, if command slice are sorted or not.
	commandsAreSorted bool
	// commandCalledAs is the name or alias value used to call this command.
	commandCalledAs struct {
		name   string
		called bool
	}

	// args is actual args parsed from flags.
	args []string
	// flagErrorBuf contains all error messages from pflag.
	flagErrorBuf *bytes.Buffer
	// flags is full set of flags.
	flags *flag.FlagSet
	// pflags contains persistent flags.
	pflags *flag.FlagSet
	// lflags contains local flags.
	lflags *flag.FlagSet
	// iflags contains inherited flags.
	iflags *flag.FlagSet
	// parentsPflags is all persistent flags of cmd's parents.
	parentsPflags *flag.FlagSet
	// globNormFunc is the global normalization function
	// that we can use on every pflag set and children commands
	globNormFunc func(f *flag.FlagSet, name string) flag.NormalizedName

	// output is an output writer defined by user.
	output io.Writer
	// versionTemplate is the version template defined by user.
	versionTemplate string
	// UsageTmpl is usage template defined by user.

	// helpTemplate is help template defined by user.
	helpTemplate string
	// helpFunc is help func defined by user.
	helpFunc func(*Command, []string)
	// helpCommand is command with usage 'help'. If it's not defined by user,
	// cobra uses default help command.
	helpCommand *Command
}

// SetArgs sets arguments for the command. It is set to os.Args[1:] by default, if desired, can be overridden
// particularly useful when testing.
func (c *Command) SetArgs(a []string) {
	c.args = a
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (c *Command) SetOutput(output io.Writer) {
	c.output = output
}

// SetGlobalNormalizationFunc sets a normalization function to all flag sets and also to child commands.
// The user should not have a cyclic dependency on commands.
func (c *Command) setGlobalNormalizationFunc(n func(f *flag.FlagSet, name string) flag.NormalizedName) {
	c.Flags().SetNormalizeFunc(n)
	c.PersistentFlags().SetNormalizeFunc(n)
	c.globNormFunc = n

	for _, command := range c.commands {
		command.setGlobalNormalizationFunc(n)
	}
}

// OutOrStdout returns output to stdout.
func (c *Command) OutOrStdout() io.Writer {
	return c.GetOut(os.Stdout)
}

// OutOrStderr returns output to stderr
func (c *Command) OutOrStderr() io.Writer {
	return c.GetOut(os.Stderr)
}

func (c *Command) GetOut(def io.Writer) io.Writer {
	if c.output != nil {
		return c.output
	}
	if c.hasParent() {
		return c.parent.GetOut(def)
	}
	return def
}

// Find the target command given the args and command tree
// Meant to be run on the highest node. Only searches down.
func (c *Command) find(args []string) (*Command, []string, error) {
	var innerfind func(*Command, []string) (*Command, []string)

	innerfind = func(c *Command, innerArgs []string) (*Command, []string) {
		argsWOflags := stripFlags(innerArgs, c)
		if len(argsWOflags) == 0 {
			return c, innerArgs
		}
		nextSubCmd := argsWOflags[0]

		cmd := c.findNext(nextSubCmd)
		if cmd != nil {
			return innerfind(cmd, argsMinusFirstX(innerArgs, nextSubCmd))
		}
		return c, innerArgs
	}

	commandFound, a := innerfind(c, args)
	if commandFound.Args == nil {
		return commandFound, a, legacyArgs(commandFound, stripFlags(a, commandFound))
	}
	return commandFound, a, nil
}

func (c *Command) findNext(next string) *Command {
	matches := make([]*Command, 0)
	for _, cmd := range c.commands {
		if cmd.Name() == next || cmd.hasAlias(next) {
			cmd.commandCalledAs.name = next
			return cmd
		}
		if EnablePrefixMatching && cmd.hasNameOrAliasPrefix(next) {
			matches = append(matches, cmd)
		}
	}

	if len(matches) == 1 {
		return matches[0]
	}

	return nil
}

// Traverse the command tree to find the command, and parse args for
// each parent.
func (c *Command) traverse(args []string) (*Command, []string, error) {
	flags := []string{}
	inFlag := false

	for i, arg := range args {
		switch {
		// A long flag with a space separated value
		case strings.HasPrefix(arg, "--") && !strings.Contains(arg, "="):
			// TODO: this isn't quite right, we should really check ahead for 'true' or 'false'
			inFlag = !hasNoOptDefVal(arg[2:], c.Flags())
			flags = append(flags, arg)
			continue
		// A short flag with a space separated value
		case strings.HasPrefix(arg, "-") && !strings.Contains(arg, "=") && len(arg) == 2 && !shortHasNoOptDefVal(arg[1:], c.Flags()):
			inFlag = true
			flags = append(flags, arg)
			continue
		// The value for a flag
		case inFlag:
			inFlag = false
			flags = append(flags, arg)
			continue
		// A flag without a value, or with an `=` separated value
		case isFlagArg(arg):
			flags = append(flags, arg)
			continue
		}

		cmd := c.findNext(arg)
		if cmd == nil {
			return c, args, nil
		}

		if err := c.parseFlags(flags); err != nil {
			return nil, args, err
		}
		return cmd.traverse(args[i+1:])
	}
	return c, args, nil
}

// VisitParents visits all parents of the command and invokes fn on each parent.
func (c *Command) VisitParents(fn func(*Command)) {
	if c.hasParent() {
		fn(c.Parent())
		c.Parent().VisitParents(fn)
	}
}

// Root finds root command.
func (c *Command) Root() *Command {
	if c.hasParent() {
		return c.Parent().Root()
	}
	return c
}

func (c *Command) execute(a []string) (err error) {
	if c == nil {
		return fmt.Errorf("Called Execute() on a nil Command")
	}
	if c.Env != nil {
		for k, v := range c.Env {
			viper.Set(k, v)
			viper.BindEnv(k, v)
		}
	}

	// initialize help and version flag at the last point possible to allow for user
	// overriding
	c.initDefaultHelpFlag()
	c.initDefaultVersionFlag()

	err = c.parseFlags(a)
	if err != nil {
		return err
	}

	helpVal, err := c.Flags().GetBool("help")
	if err != nil {
		// should be impossible to get here as we always declare a help
		// flag in InitDefaultHelpFlag()
		c.Println("\"help\" flag declared as non-bool. Please correct your code")
		return err
	}

	if helpVal {
		return flag.ErrHelp
	}

	// for back-compat, only add version flag behavior if version is defined
	if c.Version != "" {
		versionVal, err := c.Flags().GetBool("version")
		if err != nil {
			c.Println("\"version\" flag declared as non-bool. Please correct your code")
			return err
		}
		if versionVal {
			err := tmpl(c.OutOrStdout(), c.VersionTemplate(), c)
			if err != nil {
				c.Println(err)
			}
			return err
		}
	}

	if !c.Runnable() {
		return flag.ErrHelp
	}

	c.preRun()

	argWoFlags := c.Flags().Args()
	if c.DisableFlagParsing {
		argWoFlags = a
	}

	if err := c.validateArgs(argWoFlags); err != nil {
		return err
	}

	if c.PreRun != nil {
		c.PreRun(c, argWoFlags)
	}

	if c.Run != nil {
		c.Run(c, argWoFlags)
	}
	if c.PostRun != nil {
		c.PostRun(c, argWoFlags)
	}

	return nil
}

func (c *Command) preRun() {
	for _, x := range initializers {
		x()
	}
}

func (c *Command) Execute() error {
	_, err := c.executeC()
	if err != nil {
		return err
	} else {
		var srv *http.Server
		if c.Router != nil {
			if port := os.Getenv("MAMBA_PORT"); port != "" {
				srv = &http.Server{
					Handler: c.Router,
					Addr:    port,
					// Good practice: enforce timeouts for servers you create!
					WriteTimeout: 15 * time.Second,
					ReadTimeout:  15 * time.Second,
				}
			} else {
				srv = &http.Server{
					Handler: c.Router,
					Addr:    "127.0.0.1:8080",
					// Good practice: enforce timeouts for servers you create!
					WriteTimeout: 15 * time.Second,
					ReadTimeout:  15 * time.Second,
				}
			}
			return srv.ListenAndServe()
		}
	}
	return nil
}

// ExecuteC executes the command.
func (c *Command) executeC() (cmd *Command, err error) {
	// Regardless of what command execute is called on, run on Root only
	if c.hasParent() {
		return c.Root().executeC()
	}

	// windows hook
	if preExecHookFn != nil {
		preExecHookFn(c)
	}

	// initialize help as the last point possible to allow for user
	// overriding
	c.initDefaultHelpCmd()

	var args []string

	// Workaround FAIL with "go test -v" or "cobra.test -test.v", see #155
	if c.args == nil && filepath.Base(os.Args[0]) != "cobra.test" {
		args = os.Args[1:]
	} else {
		args = c.args
	}

	var flags []string
	cmd, flags, err = c.traverse(args)
	if err != nil {
		// If found parse to a subcommand and then failed, talk about the subcommand
		if cmd != nil {
			c = cmd
		}
		c.Println("Error:", err.Error())
		c.Printf("Run '%v --help' for usage.\n", c.CommandPath())

		return c, err
	}

	cmd.commandCalledAs.called = true
	if cmd.commandCalledAs.name == "" {
		cmd.commandCalledAs.name = cmd.Name()
	}

	err = cmd.execute(flags)
	if err != nil {
		// Always show help if requested, even if SilenceErrors is in
		// effect
		if err == flag.ErrHelp {
			cmd.HelpFunc()(cmd, args)
			return cmd, nil
		}

		c.Println("Error:", err.Error())
		c.Println(cmd.UsageString())
	}
	return cmd, err
}

// ResetCommands delete parent, subcommand and help command from c.
func (c *Command) ResetCommands() {
	c.parent = nil
	c.commands = nil
	c.helpCommand = nil
	c.parentsPflags = nil
}

// Commands returns a sorted slice of child commands.
func (c *Command) Commands() []*Command {
	// do not sort commands if it already sorted or sorting was disabled
	if EnableCommandSorting && !c.commandsAreSorted {
		sort.Sort(commandSorterByName(c.commands))
		c.commandsAreSorted = true
	}
	return c.commands
}

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(cmds ...*Command) {
	for i, x := range cmds {
		if cmds[i] == c {
			panic("Command can't be a child of itself")
		}
		cmds[i].parent = c
		// update max lengths
		usageLen := len(x.Use)
		if usageLen > c.commandsMaxUseLen {
			c.commandsMaxUseLen = usageLen
		}
		commandPathLen := len(x.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen {
			c.commandsMaxCommandPathLen = commandPathLen
		}
		nameLen := len(x.Name())
		if nameLen > c.commandsMaxNameLen {
			c.commandsMaxNameLen = nameLen
		}
		// If global normalization function exists, update all children
		if c.globNormFunc != nil {
			x.setGlobalNormalizationFunc(c.globNormFunc)
		}
		c.commands = append(c.commands, x)
		c.commandsAreSorted = false
		if x.Router != nil {
			c.Router.Handle("/"+x.Name(), x.Router)
		}
	}
}

// RemoveCommand removes one or more commands from a parent command.
func (c *Command) RemoveCommand(cmds ...*Command) {
	commands := []*Command{}
main:
	for _, command := range c.commands {
		for _, cmd := range cmds {
			if command == cmd {
				command.parent = nil
				continue main
			}
		}
		commands = append(commands, command)
	}
	c.commands = commands
	// recompute all lengths
	c.commandsMaxUseLen = 0
	c.commandsMaxCommandPathLen = 0
	c.commandsMaxNameLen = 0
	for _, command := range c.commands {
		usageLen := len(command.Use)
		if usageLen > c.commandsMaxUseLen {
			c.commandsMaxUseLen = usageLen
		}
		commandPathLen := len(command.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen {
			c.commandsMaxCommandPathLen = commandPathLen
		}
		nameLen := len(command.Name())
		if nameLen > c.commandsMaxNameLen {
			c.commandsMaxNameLen = nameLen
		}
	}
}

// CommandPath returns the full path to this command.
func (c *Command) CommandPath() string {
	if c.hasParent() {
		return c.Parent().CommandPath() + " " + c.Name()
	}
	return c.Name()
}

// UseLine puts out the full usage for a given command (including parents).
func (c *Command) UseLine() string {
	var useline string
	if c.hasParent() {
		useline = c.parent.CommandPath() + " " + c.Use
	} else {
		useline = c.Use
	}

	if c.hasAvailableFlags() && !strings.Contains(useline, "[flags]") {
		useline += " [flags]"
	}
	return useline
}

// Name returns the command's name: the first word in the use line.
func (c *Command) Name() string {
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// HasAlias determines if a given string is an alias of the command.
func (c *Command) hasAlias(s string) bool {
	for _, a := range c.Aliases {
		if a == s {
			return true
		}
	}
	return false
}

// CalledAs returns the command name or alias that was used to invoke
// this command or an empty string if the command has not been called.
func (c *Command) calledAs() string {
	if c.commandCalledAs.called {
		return c.commandCalledAs.name
	}
	return ""
}

// hasNameOrAliasPrefix returns true if the Name or any of aliases start
// with prefix
func (c *Command) hasNameOrAliasPrefix(prefix string) bool {
	if strings.HasPrefix(c.Name(), prefix) {
		c.commandCalledAs.name = c.Name()
		return true
	}
	for _, alias := range c.Aliases {
		if strings.HasPrefix(alias, prefix) {
			c.commandCalledAs.name = alias
			return true
		}
	}
	return false
}

// NameAndAliases returns a list of the command name and all aliases
func (c *Command) nameAndAliases() string {
	return strings.Join(append([]string{c.Name()}, c.Aliases...), ", ")
}

// Runnable determines if the command is itself runnable.
func (c *Command) Runnable() bool {
	return c.Run != nil || c.Router != nil
}

// HasSubCommands determines if the command has children commands.
func (c *Command) hasSubCommands() bool {
	return len(c.commands) > 0
}

// IsAvailableCommand determines if a command is available as a non-help command
// (this includes all non deprecated/hidden commands).
func (c *Command) isAvailableCommand() bool {
	if c.Hidden {
		return false
	}

	if c.hasParent() && c.Parent().helpCommand == c {
		return false
	}

	if c.Runnable() || c.hasAvailableSubCommands() {
		return true
	}

	return false
}

func (c *Command) isAdditionalHelpTopicCommand() bool {
	// if a command is runnable, deprecated, or hidden it is not a 'help' command
	if c.Runnable() || c.Hidden {
		return false
	}

	// if any non-help sub commands are found, the command is not a 'help' command
	for _, sub := range c.commands {
		if !sub.isAdditionalHelpTopicCommand() {
			return false
		}
	}

	// the command either has no sub commands, or no non-help sub commands
	return true
}

func (c *Command) hasHelpSubCommands() bool {
	// return true on the first found available 'help' sub command
	for _, sub := range c.commands {
		if sub.isAdditionalHelpTopicCommand() {
			return true
		}
	}

	// the command either has no sub commands, or no available 'help' sub commands
	return false
}

func (c *Command) hasAvailableSubCommands() bool {
	// return true on the first found available (non deprecated/help/hidden)
	// sub command
	for _, sub := range c.commands {
		if sub.isAvailableCommand() {
			return true
		}
	}

	return false
}

func (c *Command) hasParent() bool {
	return c.parent != nil
}