package config

import "github.com/spf13/viper"

var (
	Config = viper.New()
	Tools  = viper.New()
	Auth   = viper.New()
)

func init() {
	{
		Config.AutomaticEnv()
		Tools.AutomaticEnv()
		Auth.AutomaticEnv()
	}
	{
		Config.AllowEmptyEnv(true)
		Tools.AllowEmptyEnv(true)
		Auth.AllowEmptyEnv(true)
		{
			Config.AddConfigPath("config")
			Tools.AddConfigPath("config")
			Auth.AddConfigPath("config")

			Config.SetConfigType("yaml")
			Tools.AddConfigPath("yaml")
			Auth.AddConfigPath("yaml")
		}
		{
			Config.SetConfigName("config")
			Tools.SetConfigName("tools")
			Auth.SetConfigName("auth")
		}
		Config.SafeWriteConfig()
		Tools.SafeWriteConfig()
		Auth.SafeWriteConfig()
	}
	{
		Config.ReadInConfig()
		Tools.ReadInConfig()
		Auth.ReadInConfig()
	}

}
