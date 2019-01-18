package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"strings"
	"github.com/fatih/color"
)

// FlagValueSet is an interface that users can implement
// to bind a set of flags to config.
type FlagValueSet interface {
	VisitAll(fn func(FlagValue))
}

// FlagValue is an interface that users can implement
// to bind different flags to config.
type FlagValue interface {
	HasChanged() bool
	Name() string
	ValueString() string
	ValueType() string
}

// pflagValueSet is a wrapper around *pflag.ValueSet
// that implements FlagValueSet.
type pflagValueSet struct {
	flags *pflag.FlagSet
}

// VisitAll iterates over all *pflag.Flag inside the *pflag.FlagSet.
func (p pflagValueSet) VisitAll(fn func(flag FlagValue)) {
	p.flags.VisitAll(func(flag *pflag.Flag) {
		fn(pflagValue{flag})
	})
}

// pflagValue is a wrapper aroung *pflag.flag
// that implements FlagValue
type pflagValue struct {
	flag *pflag.Flag
}

// HasChanges returns whether the flag has changes or not.
func (p pflagValue) HasChanged() bool {
	return p.flag.Changed
}

// Name returns the name of the flag.
func (p pflagValue) Name() string {
	return p.flag.Name
}

// ValueString returns the value of the flag as a string.
func (p pflagValue) ValueString() string {
	return p.flag.Value.String()
}

// ValueType returns the type of the flag as a string.
func (p pflagValue) ValueType() string {
	return p.flag.Value.Type()
}

func (v *config) BindPFlags(flags *pflag.FlagSet) error {
	return v.BindFlagValues(pflagValueSet{flags})
}

func (v *config) BindPFlag(key string, flag *pflag.Flag) error {
	return v.BindFlagValue(key, pflagValue{flag})
}

func (v *config) BindFlagValues(flags FlagValueSet) (err error) {
	flags.VisitAll(func(flag FlagValue) {
		if err = v.BindFlagValue(flag.Name(), flag); err != nil {
			return
		}
	})
	return nil
}

func (v *config) BindFlagValue(key string, flag FlagValue) error {
	if flag == nil {
		return fmt.Errorf(color.RedString("flag for %q is nil", key))
	}
	v.pflags[strings.ToLower(key)] = flag
	return nil
}
