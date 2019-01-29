package config

import (
	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"runtime"
)

var (
	homeDir, _ = homedir.Dir()
	Viper      = viper.New()
	sys        = &afero.Afero{
		Fs: afero.NewOsFs(),
	}
)

var Configuration = &Config{}

// Config contains service configuration
type Config struct {
	Project struct {
		Name        string `mapstructure:"project.name"`
		Description string `mapstructure:"project.description"`
		Github      string `mapstructure:"project.github"`
		Project     string `mapstructure:"project.project"`
		Bin         string `mapstructure:"project.bin"`
		GitInit     bool   `mapstructure:"project.gitinit"`
		Contract    bool   `mapstructure:"project.contract"`
	}
	GKE struct {
		Enabled bool   `mapstructure:"gke.enabled"`
		Project string `mapstructure:"gke.project"`
		Zone    string `mapstructure:"gke.zone"`
		Cluster string `mapstructure:"gke.cluster"`
	}
	Storage struct {
		Enabled  bool `mapstructure:"storage.enabled"`
		Postgres bool `mapstructure:"storage.posgres"`
		MySQL    bool `mapstructure:"storage.mysql"`
		Config   struct {
			Driver      string `mapstructure:"storage.comfig.driver"`
			Host        string `mapstructure:"storage.comfig.driver"`
			Port        int    `mapstructure:"storage.comfig.driver"`
			Name        string `mapstructure:"storage.comfig.driver"`
			Username    string `mapstructure:"storage.comfig.driver"`
			Password    string `mapstructure:"storage.comfig.driver"`
			Connections struct {
				Max  int `mapstructure:"storage.comfig.connections.max"`
				Idle int `mapstructure:"storage.comfig.connections.idle"`
			}
		}
	}
	API struct {
		Enabled bool `mapstructure:"api.enabled"`
		GRPC    bool `mapstructure:"api.grpc"`
		Gateway bool `mapstructure:"api.gateway"`
		Config  struct {
			Port    int `mapstructure:"api.config.port"`
			Gateway struct {
				Port int `mapstructure:"api.config.port"`
			}
		}
	}
	Directories struct {
		Templates string `mapstructure:"directories.templates"`
		Service   string `mapstructure:"directories.service"`
	}
	Docker struct {
		Endpoint  string `mapstructure:"docker.endpoint"`
		Dockerhub string `mapstructure:"docker.dockerhub"`
		Container struct {
			Image      string   `mapstructure:"docker.container.image"`
			Container  string   `mapstructure:"docker.container"`
			Args       []string `mapstructure:"docker.container.dependencies"`
			Env        []string `mapstructure:"docker.container.env"`
			Entrypoint []string `mapstructure:"docker.container.entrypoint"`
			Modules    bool     `mapstructure:"docker.container.modules"`
		}
	}
}

func init() {
	Viper.SetConfigName(FileName)
	Viper.AllowEmptyEnv(true)
	Viper.SetConfigType("json")
	Viper.AddConfigPath(homeDir)
	Viper.AddConfigPath(".")
	Viper.AutomaticEnv()
	zap.LogE("Writing config", viper.WriteConfig())
	{
		Viper.SetDefault("project.name", "default")
		Viper.SetDefault("project.description", "default")
		Viper.SetDefault("project.github", "default")
		Viper.SetDefault("project.project", "default")
		Viper.SetDefault("project.bin", "default")
		Viper.SetDefault("project.gitinit", "default")
		Viper.SetDefault("project.contract", "default")

	}
	{
		Viper.SetDefault("gke.enabled", "default")
		Viper.SetDefault("gke.project", "default")
		Viper.SetDefault("gke.zone", "default")
		Viper.SetDefault("gke.cluster", "default")
	}
	{
		Viper.SetDefault("storage.enabled", true)
		Viper.SetDefault("storage.postgress", true)
		Viper.SetDefault("storage.mysql", false)
		Viper.SetDefault("storage.config.driver", "default")
		Viper.SetDefault("storage.config.host", "default")
		Viper.SetDefault("storage.config.port", "default")
		Viper.SetDefault("storage.config.name", "default")
		Viper.SetDefault("storage.config.username", "default")
		Viper.SetDefault("storage.config.password", "default")
		Viper.SetDefault("storage.config.connections.max", 10)
		Viper.SetDefault("storage.config.connections.idle", 1)
	}
	{
		Viper.SetDefault("api.enabled", true)
		Viper.SetDefault("api.grpc", true)
		Viper.SetDefault("api.gateway", true)
		Viper.SetDefault("api.config.port", "8080")
	}
	{
		Viper.SetDefault("directories.templates", "")
		Viper.SetDefault("directories.service", "")
	}
	{
		Viper.SetDefault("docker.endpoint", "unix:///var/run/docker.sock")
		Viper.SetDefault("docker.dockerhub", "colemanword")
		Viper.SetDefault("docker.container.image", "golang:1.11")
		Viper.SetDefault("docker.container.modules", true)
	}

	{
		Viper.SetDefault("os.runtime.goarch", runtime.GOARCH)
		Viper.SetDefault("os.runtime.compiler", runtime.Compiler)
		Viper.SetDefault("os.runtime.version", runtime.Version())
		Viper.SetDefault("os.runtime.goos", runtime.GOOS)
	}
	{
		Viper.SetDefault("config.name", FileName)
	}
	zap.LogE("updating config", Viper.WriteConfig())
	zap.LogE("Reading config", Viper.ReadInConfig())
	zap.LogF("unmarshaling config", Viper.Unmarshal(Configuration))
	zap.Debug("Current config file-->", "config", Viper.ConfigFileUsed())
}

func GetConfig() *Config {
	return Configuration
}

func Annotate(v *viper.Viper) map[string]string {
	settings := v.AllSettings()
	an := make(map[string]string)
	for k, v := range settings {
		if t, ok := v.(string); ok == true {
			an[k] = t
		}
	}
	return an
}