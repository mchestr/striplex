package config

import (
	"strings"

	"github.com/kokizzu/gotro/L"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	setDefaults()
	config.SetConfigType("yaml")
	config.SetConfigName("default")
	config.AddConfigPath("config/")
	config.SetEnvPrefix("striplex")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.SetConfigType("env")
	config.AutomaticEnv()
	err = config.ReadInConfig()
	L.PanicIf(err, "error on parsing default environment configuration file", err)

	// Merge the environment configuration file
	envFileConfig := viper.New()
	envFileConfig.SetConfigType("yaml")
	envFileConfig.AddConfigPath("config/")
	envFileConfig.SetConfigName(env)
	err = envFileConfig.ReadInConfig()
	L.PanicIf(err, "error on parsing environment configuration file", err)

	err = config.MergeConfigMap(envFileConfig.AllSettings())
	L.PanicIf(err, "error on merging environment configuration file", err)
}

func GetConfig() *viper.Viper {
	return config
}

func setDefaults() {
	config.SetDefault("server.mode", "release")
	config.SetDefault("server.address", ":8080")
}
