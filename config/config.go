package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

var Config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error

	// Merge the environment configuration file
	defaultFileConfig := viper.New()
	defaultFileConfig.AddConfigPath("config/")
	defaultFileConfig.SetConfigName("default")
	if err = defaultFileConfig.ReadInConfig(); err != nil {
		slog.Error("error on parsing default configuration file", "error", err)
		panic(err)
	}

	// Merge the environment configuration file
	envFileConfig := viper.New()
	envFileConfig.AddConfigPath("config/")
	envFileConfig.SetConfigName(env)
	if err = envFileConfig.ReadInConfig(); err != nil {
		slog.Error("error on parsing environment configuration file", "env", env, "error", err)
		panic(err)
	}

	if err = defaultFileConfig.MergeConfigMap(envFileConfig.AllSettings()); err != nil {
		slog.Error("error on merging configuration file", "env", env, "error", err)
		panic(err)
	}

	Config = viper.New()
	if err = Config.MergeConfigMap(defaultFileConfig.AllSettings()); err != nil {
		slog.Error("error on merging default configuration", "error", err)
		panic(err)
	}
	Config.SetConfigType("env")
	Config.SetEnvPrefix("striplex")
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	Config.AutomaticEnv()
	setDefaults()
}

func setDefaults() {
	Config.SetDefault("server.mode", "release")
	Config.SetDefault("server.address", ":8080")
}
