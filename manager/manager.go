package manager

import (
	"github.com/go-kit/kit/log"
	"github.com/gofunct/mamba/input"
	"github.com/gofunct/mamba/logging"
	"github.com/gofunct/mamba/walker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Interface interface {
	Set(s interface{})
	GetString(s string) string
	Unmarshal(interface{}, []byte) error
	log.Logger

}

type Manager struct {
	L *logging.CtxLogger
	Q *input.UI
	Flags *pflag.FlagSet
	W *walker.Walker
}

func (m *Manager) Set(s string, k interface{}) {
	viper.Set(s, k)
}

func (m *Manager) GetString(s string) string {
	if res := viper.GetString(s); res != "" {
		return res
	}

	if res := os.Getenv(strings.ToUpper(s)); res != "" {
		viper.Set(s, res)
		return res
	}
	if res := m.Q.Enquire("Please initialize variable:", s); res != "" {
		viper.Set(s, res)
		return res
	}
	return ""
}

func NewManager(set *pflag.FlagSet) *Manager {
	return &Manager{
		L: logging.NewLogCtx(logrus.New()),
		Q: input.DefaultUI(),
		W: walker.NewWalker(),
		Flags: set,
	}
}
