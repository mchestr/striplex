package config

import (
	"strings"

	"github.com/kokizzu/gotro/L"
	"github.com/spf13/viper"
)

var Config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	Config = viper.New()
	setDefaults()
	Config.SetConfigType("yaml")
	Config.SetConfigName("default")
	Config.AddConfigPath("config/")
	Config.SetEnvPrefix("striplex")
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	Config.SetConfigType("env")
	Config.AutomaticEnv()
	err = Config.ReadInConfig()
	L.PanicIf(err, "error on parsing default environment configuration file", err)

	// Merge the environment configuration file
	envFileConfig := viper.New()
	envFileConfig.SetConfigType("yaml")
	envFileConfig.AddConfigPath("config/")
	envFileConfig.SetConfigName(env)
	err = envFileConfig.ReadInConfig()
	L.PanicIf(err, "error on parsing environment configuration file", err)

	err = Config.MergeConfigMap(envFileConfig.AllSettings())
	L.PanicIf(err, "error on merging environment configuration file", err)
}

func setDefaults() {
	Config.SetDefault("server.mode", "release")
	Config.SetDefault("server.address", ":8080")
}
