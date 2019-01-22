package manager

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gofunct/mamba/manager/input"
	"github.com/gofunct/mamba/logging"
	"github.com/gofunct/mamba/walker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
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
	Requirements []string
}

func (m *Manager) Set(s string, k interface{}) {
	viper.Set(s, k)
	_ = os.Setenv(s, k.(string))
	_ = m.Flags.Set(s, k.(string))
}

func (m *Manager) GetString(s string) string {
	if res := viper.GetString(s); res != "" {
		return res
	}
	if res, err := m.Flags.GetString(s); res != "" && err != nil {
		return res
	}

	if res := os.Getenv(strings.ToUpper(s)); res != "" {
		viper.Set(s, res)
		return res
	}
	if res := m.Q.Enquire("Please set a value for variable: "+s, s); res != "" {
		viper.Set(s, res)
		return res
	}
	return ""
}

func NewManager(vars []string) *Manager {
	m := &Manager{
		L: logging.NewLogCtx(logrus.New()),
		Q: input.DefaultUI(),
		W: walker.NewWalker(),
		Requirements: vars,
	}
	return m
}

func (m *Manager) SyncRequirements()  {
	for _, s := range m.Requirements {
		if res := viper.GetString(s); res != "" {
			if err := viper.BindEnv(s, res); err != nil {
				m.L.Warn("failed to bind env", s, res)
			}
		}
		if res, err := m.Flags.GetString(s); res != "" && err != nil {
			viper.Set(s, res)
			if err := viper.BindEnv(s, res); err != nil {
				m.L.Warn("failed to bind env", s, res)
			}
		}

		if res := os.Getenv(strings.ToUpper(s)); res != "" {
			viper.Set(s, res)
			if err := os.Setenv(strings.ToUpper(s), res); err != nil {
				m.L.Warn("failed to set env", strings.ToUpper(s), res)
			}
		}
		if res := m.Q.Enquire("Please set a value for variable: "+s, s); res != "" {
			viper.Set(s, res)
			if err := viper.BindEnv(s, res); err != nil {
				m.L.Warn("failed to bind env", s, res)
			}
		}
	}
	m.Write()
	m.L.Debug("all requirements have been synced")
}

func (m *Manager) AddFlagSet(set *pflag.FlagSet) {
	m.Flags = set
	if err := viper.BindPFlags(m.Flags); err != nil {
		m.L.Warn("failed to bind to pflags", err.Error())
	}
}

func (m *Manager) GetRequirements() map[string]string {
	set := make(map[string]string)
	for _, v := range m.Requirements {
		set[v] = m.GetString(v)
	}
	return set
}

func (m *Manager) Debug() {
	fmt.Println("Config:")
	viper.Debug()
	fmt.Println("Requirements: ", m.Requirements)
	for _, v := range m.Q.Queries {
		fmt.Println("Question: ",v.Q)
		fmt.Println("Tag: ",v.Tag)
		fmt.Println("Name: ",v.Opts.Name)
		fmt.Println("Required: ", v.Opts.Required)
		fmt.Println("Default: ", v.Opts.Default)
	}
}

func (m *Manager) Write() {
	if res, err := m.Q.Select("Write current config to disc?", []string{"yes", "no"}, &input.Options{});
	err == nil && strings.Contains(res, "y") || err == nil && strings.Contains(res, "Y") || err == nil && strings.Contains(res, "yes") {
		if file := m.Q.Enquire("Please provide a path to config file", "file"); file != "" {
			if err := viper.WriteConfigAs(file); err != nil {
				m.L.Warn("failed to write config", err.Error())
			}
			m.L.Debug("Updated config file successfully")
		}
	}
}

func (m *Manager) Unmarshal(i interface{}) {
	if err := viper.Unmarshal(i); err != nil {
		m.L.Warn("failed to unmarshal", err.Error())
	}
}

func (m *Manager) AllSettings() map[string]interface{} {
	if err := viper.MergeInConfig(); err != nil {
		m.L.Warn("failed to merge in config", err.Error())
	}
	return viper.AllSettings()
}


func (m *Manager) WriteFile(f string, d []byte)  error {
	return ioutil.WriteFile(f, d, 0755)
}

func (m *Manager) ReadFile(f string)  ([]byte, error) {
	return ioutil.ReadFile(f)
}

func (m *Manager) ReadStdIn()  ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}
func (m *Manager) ReadReader(reader io.Reader)  ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func (m *Manager) ReadDir(f string)  ([]os.FileInfo, error) {
	return ioutil.ReadDir(f)
}