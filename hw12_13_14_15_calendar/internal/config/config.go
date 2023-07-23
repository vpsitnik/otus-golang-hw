package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger *LoggerConf
}

type LoggerConf struct {
	Level string
}

func NewConfig(configFile string) *Config {
	viper.SetDefault("logger.level", "INFO")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &Config{
		Logger: &LoggerConf{
			Level: viper.GetString("logger.level"),
		},
	}
}

// TODO
