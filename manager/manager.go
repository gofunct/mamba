package manager

import (
	"fmt"
	"github.com/gofunct/mamba/manager/input"
	"github.com/gofunct/mamba/manager/logging"
	"github.com/gofunct/mamba/manager/walker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Manager struct {
	l            *logging.CtxLogger
	Q            *input.UI
	flags        *pflag.FlagSet
	W            *walker.Walker
	Requirements []string
	usg          string
}

func (m *Manager) Set(s string, k interface{}, req bool) {
	viper.SetDefault(s, k)
	_ = os.Setenv(s, k.(string))
	_ = m.flags.Set(s, k.(string))
	if req {
		m.AddRequirement(s)
	}
}

func (m *Manager) GetString(s string) string {

	if res := viper.GetString(s); res != "" {
		return res
	}
	if res, err := m.flags.GetString(s); res != "" && err != nil {
		return res
	}

	if res := os.Getenv(strings.ToUpper(s)); res != "" {
		return res
	}
	q := &input.Query{
		Q: "Please provide a value for: " + s,
		Opts: &input.Options{
			Key:      s,
			Loop:     true,
			Required: true,
		},
	}
	if res, err := m.Q.Ask(q); res != "" && err == nil {
		return res
	} else {
		m.Fatalf("%s failed to retrieve value for: %s\n%s", s, err.Error())
	}
	return ""
}

func NewManager() *Manager {
	m := &Manager{
		l: logging.NewLogCtx(logrus.New()),
		Q: input.DefaultUI(),
		W: walker.NewWalker(),
	}
	return m
}

func (m *Manager) SyncRequirements() {
	for _, s := range m.Requirements {
		if res := viper.GetString(s); res != "" {
			if err := viper.BindEnv(s, res); err != nil {
				m.l.Warn("failed to bind env", s, res)
			}
		}
		if res, err := m.flags.GetString(s); res != "" && err != nil {
			viper.Set(s, res)
			if err := viper.BindEnv(s, res); err != nil {
				m.l.Warn("failed to bind env", s, res)
			}
		}

		if res := os.Getenv(strings.ToUpper(s)); res != "" {
			viper.Set(s, res)
			if err := os.Setenv(strings.ToUpper(s), res); err != nil {
				m.l.Warn("failed to set env", strings.ToUpper(s), res)
			}
		}
		needed := m.Q.SingleQuery("Please set a value for variable: "+s, &input.Options{
			Key:      s,
			Loop:     true,
			Required: true,
		})
		if res, err := m.Q.Ask(needed); err == nil && res != "" {
			viper.Set(s, res)
			if err := viper.BindEnv(s, res); err != nil {
				m.l.Warn("failed to bind env", s, res)
			}
		}
	}
	m.Write()
	m.l.Debug("all requirements have been synced")
}

func (m *Manager) AddFlagSet(set *pflag.FlagSet) {
	m.flags = set
	if err := viper.BindPFlags(m.flags); err != nil {
		m.l.Warn("failed to bind to pflags", err.Error())
	}
}

func (m *Manager) GetRequirements() map[string]string {
	set := make(map[string]string)
	for _, v := range m.Requirements {
		set[v] = m.GetString(v)
	}
	return set
}

func (m *Manager) Help() {
	fmt.Println("Config:")
	viper.Debug()
	fmt.Println("Requirements: ", m.Requirements)
	for _, v := range m.Q.Queries {
		fmt.Println("Question: ", v.Q)
		fmt.Println("Key: ", v.Opts.Key)
		fmt.Println("Options: ", v.Opts.Options)
		fmt.Println("Required: ", v.Opts.Required)
		fmt.Println("Default: ", v.Opts.Default)
		fmt.Println("Loop: ", v.Opts.Loop)

	}
}

func (m *Manager) Write() {
	sel := m.Q.SingleQuery("Write current config to disc?", &input.Options{
		Key:      "config.write",
		Options:  []string{"true", "false"},
		Default:  "false",
		Loop:     true,
		Required: true,
	})
	cfg := m.Q.SingleQuery("Please provide a path to config file", &input.Options{
		Key:      "config.path",
		Options:  []string{"y", "n"},
		Default:  "n",
		Loop:     true,
		Required: true,
	})

	if res, err := m.Q.Select(sel); res == "true" && err == nil {
		if file, err := m.Q.Ask(cfg); err == nil && file != "" {
			if err := viper.WriteConfigAs(file); err != nil {
				m.l.Warn("failed to write config", err.Error())
			}
			m.l.Debug("Updated config file successfully")
		}
	}
}

func (m *Manager) Unmarshal(i interface{}) {
	if err := viper.Unmarshal(i); err != nil {
		m.l.Warn("failed to unmarshal", err.Error())
	}
}

func (m *Manager) AllSettings() map[string]interface{} {
	if err := viper.MergeInConfig(); err != nil {
		m.l.Warn("failed to merge in config", err.Error())
	}
	return viper.AllSettings()
}

func (m *Manager) Usage() {
	fmt.Println(m.usg)
}

func (m *Manager) WriteFile(f string, d []byte) error {
	return ioutil.WriteFile(f, d, 0755)
}

func (m *Manager) ReadFile(f string) ([]byte, error) {
	return ioutil.ReadFile(f)
}

func (m *Manager) ReadStdIn() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}
func (m *Manager) ReadReader(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func (m *Manager) ReadDir(f string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(f)
}

func (m *Manager) BindCmd(cmd *cobra.Command) {
	m.AddFlagSet(cmd.Flags())
	m.AddFlagSet(cmd.PersistentFlags())
	m.usg = cmd.UsageString()
}

func (m *Manager) AddRequirement(s string) {
	m.Requirements = append(m.Requirements, s)
}

func (m *Manager) Warnf(f string, args ...interface{}) {
	m.l.Warnf(f, args)
}
func (m *Manager) Fatalf(f string, args ...interface{}) {
	m.l.Fatalf(f, args)
}

func (m *Manager) Debug(args ...interface{}) {
	m.l.Debug(args)
}

func (m *Manager) Log(args ...interface{}) {
	if err := m.l.Log(args); err != nil {
		m.l.Warn("failed to log context", err.Error())
	}
}

func (m *Manager) SetUsage(s string) {
	m.usg = s
}
