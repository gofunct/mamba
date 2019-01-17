package cli

import "github.com/spf13/cobra"

type CobraFunc func(*cobra.Command, []string)
