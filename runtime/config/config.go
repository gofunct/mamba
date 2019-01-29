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
		Name        string `json:"project.name"`
		Description string `json:"project.description"`
		Github      string `json:"project.github"`
		Project     string `json:"project.project"`
		Bin         string `json:"project.bin"`
		GitInit     bool   `json:"project.gitinit"`
		Contract    bool   `json:"project.contract"`
	}
	GKE struct {
		Enabled bool   `json:"gke.enabled"`
		Project string `json:"gke.project"`
		Zone    string `json:"gke.zone"`
		Cluster string `json:"gke.cluster"`
	}
	Storage struct {
		Enabled  bool `json:"storage.enabled"`
		Postgres bool `json:"storage.posgres"`
		MySQL    bool `json:"storage.mysql"`
		Config   struct {
			Driver      string `json:"storage.comfig.driver"`
			Host        string `json:"storage.comfig.driver"`
			Port        int    `json:"storage.comfig.driver"`
			Name        string `json:"storage.comfig.driver"`
			Username    string `json:"storage.comfig.driver"`
			Password    string `json:"storage.comfig.driver"`
			Connections struct {
				Max  int `json:"storage.comfig.connections.max"`
				Idle int `json:"storage.comfig.connections.idle"`
			}
		}
	}
	API struct {
		Enabled bool `json:"api.enabled"`
		GRPC    bool `json:"api.grpc"`
		Gateway bool `json:"api.gateway"`
		Config  struct {
			Port    int `json:"api.config.port"`
			Gateway struct {
				Port int `json:"api.config.port"`
			}
		}
	}
	Directories struct {
		Templates string `json:"directories.templates"`
		Service   string `json:"directories.service"`
	}
	Docker struct {
		Endpoint  string `json:"docker.endpoint"`
		Dockerhub string `json:"docker.dockerhub"`
		Container struct {
			Image      string   `json:"docker.container.image"`
			Container  string   `json:"docker.container"`
			Args       []string `json:"docker.container.dependencies"`
			Env        []string `json:"docker.container.env"`
			Entrypoint []string `json:"docker.container.entrypoint"`
			Modules    bool     `json:"docker.container.modules"`
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
	zap.LogE("Writing config", Viper.WriteConfig())
	{
		Viper.SetDefault("project.name", "default")
		Viper.SetDefault("project.description", "default")
		Viper.SetDefault("project.github", "default")
		Viper.SetDefault("project.project", "default")
		Viper.SetDefault("project.bin", "default")
		Viper.SetDefault("project.gitinit", false)
		Viper.SetDefault("project.contract", true)

	}
	{
		Viper.SetDefault("gke.enabled", true)
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
		Viper.SetDefault("storage.config.port", DefaultPostgresPort)
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
	Viper.Unmarshal(Configuration)
	return Configuration
}