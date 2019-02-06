package input

import (
	"github.com/gofunct/mamba/runtime/logging"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type ValidateFunc func(string) error

type Options struct {
	Key        string
	Options    []string
	Default    string
	Loop       bool
	Required   bool
	hasChanged bool
}

func bind(q *Query) ValidateFunc {
	return func(ans string) error {
		if ans == "" {
			return ErrEmpty
		}
		viper.Set(q.Opts.Key, ans)

		if err := os.Setenv(strings.ToUpper(q.Opts.Key), ans); err != nil {
			logging.L.Warn("failed to set env", err.Error())
		}
		return nil
	}
}
